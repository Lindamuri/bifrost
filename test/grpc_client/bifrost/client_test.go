package bifrost

import (
	"context"
	"fmt"
	nginx_ctx "github.com/tremendouscan/bifrost/pkg/resolv/V3/nginx/configuration/context"
	"github.com/tremendouscan/bifrost/pkg/resolv/V3/nginx/configuration/context/local"
	"github.com/tremendouscan/bifrost/pkg/resolv/V3/nginx/configuration/context_type"
	utilsV3 "github.com/tremendouscan/bifrost/pkg/resolv/V3/nginx/configuration/utils"
	"github.com/marmotedu/errors"
	"sync"
	"testing"
	"time"

	"google.golang.org/grpc"

	v1 "github.com/tremendouscan/bifrost/api/bifrost/v1"
	bifrost_cliv1 "github.com/tremendouscan/bifrost/pkg/client/bifrost/v1"
	healthzclient_v1 "github.com/tremendouscan/bifrost/pkg/client/grpc_health_v1"
)

func TestRun(t *testing.T) {
	err := exampleServerRun()
	if err != nil {
		t.Fatal(err.Error())
	}
}

func TestBifrostClient(t *testing.T) {
	healthCli, err := healthzclient_v1.NewClient(serverAddress(), grpc.WithInsecure())
	if err != nil {
		t.Fatalf(err.Error())
	}

	retryI := 0
	for {
		if retryI >= 10 {
			t.Fatalf("connect to web server config service timeout.")
		}
		state, err := healthCli.Check(context.Background(), "com.github.ClessLi.api.bifrostpb.WebServerConfig")
		if err != nil {
			t.Log(err.Error())
		}
		if state == healthzclient_v1.SERVING {
			t.Log("service is serving")
			break
		}
		time.Sleep(time.Second * 3)
		retryI++
	}

	client, err := bifrost_cliv1.New(serverAddress(), grpc.WithInsecure(), grpc.WithTimeout(time.Second))
	if err != nil {
		t.Fatalf(err.Error())
	}

	defer client.Close()

	servernames, err := client.WebServerConfig().GetServerNames()
	if err != nil {
		t.Fatalf(err.Error())
	}

	time.Sleep(time.Second * 10)
	metrics, err := client.WebServerStatus().Get()
	if err != nil {
		t.Fatalf("%++v", err)
	}
	t.Log(metrics)

	// normal grpc client
	/*cclient, err := grpc.Dial(serverAddress(), grpc.WithInsecure(), grpc.WithTimeout(time.Second))
	if err != nil {
		t.Fatalf(err.Error())
	}

	defer cclient.Close()
	c := pbv1.NewWebServerConfigClient(cclient)*/
	wg := new(sync.WaitGroup)
	for _, servername := range servernames {
		// normal grpc client
		/*resp, err := c.Get(context.Background(), &pbv1.ServerName{Name: servername})
		if err != nil {
			t.Fatalf(err.Error())
		}
		buf := bytes.NewBuffer(nil)
		stop := false
		for !stop {
			select {
			case <-resp.Context().Done():
				stop = true
				break
			default:
				conf, err := resp.Recv()
				if err != nil && err != io.EOF {
					t.Fatalf(err.Error())
				}
				buf.Write(conf.GetJsonData())
				if err == io.EOF {
					stop = true
					break
				}
			}
		}
		t.Logf("config:\n\n%s", buf.String())*/

		// go-kit grpc client
		conf, fingerprinter, err := client.WebServerConfig().Get(servername)
		if err != nil {
			t.Fatalf(err.Error())
		}
		t.Logf("get the config original fingerprints: %v", fingerprinter.Fingerprints())
		lines, err := conf.Main().ConfigLines(false)
		if err != nil {
			t.Fatalf(err.Error())
		}
		t.Logf("get config lines len: %d", len(lines))
		for _, line := range lines {
			t.Log(line)
		}

		statistics, err := client.WebServerStatistics().Get(servername)
		if err != nil {
			t.Fatalf(err.Error())
		}
		t.Logf("statistics %s:\n\n%+v", servername, statistics)

		logC, lwCancel, err := client.WebServerLogWatcher().Watch(&v1.WebServerLogWatchRequest{
			ServerName:          &v1.ServerName{Name: servername},
			LogName:             "access.log",
			FilteringRegexpRule: "^test.*$",
		})
		if err != nil {
			t.Fatalf(err.Error())
		}

		wg.Add(1)
		go func() {
			defer wg.Done()
			defer lwCancel()
			for {
				select {
				case <-time.After(time.Second * 10):
					return
				case line := <-logC:
					if line == nil {
						return
					}
					fmt.Print(string(line))
				}
			}
		}()
	}
	wg.Wait()
}

func BenchmarkBifrostClientGet(b *testing.B) {
	client, err := bifrost_cliv1.New(serverAddress(), grpc.WithInsecure(), grpc.WithTimeout(time.Second))
	if err != nil {
		b.Fatalf(err.Error())
	}

	defer client.Close()
	for i := 0; i < b.N; i++ {
		_, _, err := client.WebServerConfig().Get("example test")
		if err != nil {
			b.Fatalf(err.Error())
		}
	}
}

func TestBifrostClientOperation(t *testing.T) {
	client, err := bifrost_cliv1.New(serverAddress(), grpc.WithInsecure(), grpc.WithTimeout(time.Second))
	if err != nil {
		t.Fatal(err)
	}

	defer client.Close()

	servernames, err := client.WebServerConfig().GetServerNames()
	if err != nil {
		t.Fatal(err)
	}

	for _, servername := range servernames {
		conf, ofp, err := client.WebServerConfig().Get(servername)
		if err != nil {
			t.Fatal(err)
		}
		ctx, idx := conf.Main().ChildrenPosSet().
			QueryOne(nginx_ctx.NewKeyWords(context_type.TypeHttp).
				SetSkipQueryFilter(nginx_ctx.SkipDisabledCtxFilterFunc)).
			QueryAll(nginx_ctx.NewKeyWords(context_type.TypeServer).
				SetSkipQueryFilter(nginx_ctx.SkipDisabledCtxFilterFunc)).
			Filter(
				func(pos nginx_ctx.Pos) bool {
					return pos.QueryOne(nginx_ctx.NewKeyWords(context_type.TypeDirective).
						SetCascaded(false).
						SetStringMatchingValue("server_name test1.com").
						SetSkipQueryFilter(nginx_ctx.SkipDisabledCtxFilterFunc)).
						Target().Error() == nil
				},
			).
			Filter(
				func(pos nginx_ctx.Pos) bool {
					return pos.QueryOne(nginx_ctx.NewKeyWords(context_type.TypeDirective).
						SetCascaded(false).
						SetRegexpMatchingValue("^listen 80$").
						SetSkipQueryFilter(nginx_ctx.SkipDisabledCtxFilterFunc)).
						Target().Error() == nil
				},
			).
			QueryOne(nginx_ctx.NewKeyWords(context_type.TypeLocation).
				SetRegexpMatchingValue(`^/test1-location$`).
				SetSkipQueryFilter(nginx_ctx.SkipDisabledCtxFilterFunc)).
			QueryOne(nginx_ctx.NewKeyWords(context_type.TypeIf).
				SetRegexpMatchingValue(`^\(\$http_api_name != ''\)$`).
				SetSkipQueryFilter(nginx_ctx.SkipDisabledCtxFilterFunc)).
			QueryOne(nginx_ctx.NewKeyWords(context_type.TypeDirective).
				SetStringMatchingValue("proxy_pass").
				SetSkipQueryFilter(nginx_ctx.SkipDisabledCtxFilterFunc)).
			Position()
		err = ctx.Insert(local.NewContext(context_type.TypeInlineComment, fmt.Sprintf("[%s]test comments", time.Now().String())), idx+1).Error()
		if err != nil {
			t.Fatal(err)
		}
		if err := client.WebServerConfig().Update(servername, conf, utilsV3.NewConfigFingerprinter(conf.Dump()).Fingerprints()); err == nil {
			t.Fatal(errors.New("the original fingerprint was not used for updating, but the update did not return an error"))
		} else {
			t.Log("the original fingerprint was not used for updating, and the update return an error")
		}
		err = client.WebServerConfig().Update(servername, conf, ofp.Fingerprints())
		if err != nil {
			t.Fatal(err)
		}
	}
}

func TestBifrostClientExecServerBinCMD(t *testing.T) {
	client, err := bifrost_cliv1.New(serverAddress(), grpc.WithInsecure(), grpc.WithTimeout(time.Second))
	if err != nil {
		t.Fatal(err)
	}

	defer client.Close()

	servernames, err := client.WebServerConfig().GetServerNames()
	if err != nil {
		t.Fatal(err)
	}

	for _, servername := range servernames {
		// nil args
		s1, out1, err1, err := client.WebServerBinCMD().Exec(servername)
		if err != nil {
			t.Errorf("[nil args] server name: %s, error: %v", servername, err)
		}
		t.Logf("[nil args] server name: %s, exec successful: %v, result stdout: %s, result stderr: %s", servername, s1, out1, err1)

		// one arg
		s2, out2, err2, err := client.WebServerBinCMD().Exec(servername, "-t")
		if err != nil {
			t.Errorf("[one arg] server name: %s, error: %v", servername, err)
		}
		t.Logf("[one arg] server name: %s, exec successful: %v, result stdout: %s, result stderr: %s", servername, s2, out2, err2)

		// two args
		s3, out3, err3, err := client.WebServerBinCMD().Exec(servername, "-s", "reload")
		if err != nil {
			t.Errorf("[two args] server name: %s, error: %v", servername, err)
		}
		t.Logf("[two args] server name: %s, exec successful: %v, result stdout: %s, result stderr: %s", servername, s3, out3, err3)
	}
}
