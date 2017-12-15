/* package gowechat 一个简单易用的wechat封装.
 */
package gowechat

import (
	"fmt"
	"sync"

	"github.com/astaxie/beego/cache"
	"github.com/yaotian/gowechat/util"
	"github.com/yaotian/gowechat/wxcontext"
)

// Wechat struct
type Wechat struct {
	Context *wxcontext.Context
}

// NewWechat init
func NewWechat(cfg wxcontext.Config) *Wechat {
	context := new(wxcontext.Context)
	initContext(cfg, context)
	return &Wechat{context}
}

func initContext(cfg wxcontext.Config, context *wxcontext.Context) {
	if cfg.Cache == nil {
		cfg.Cache, _ = cache.NewCache("memory", `{"interval":60}`)
	}
	context.Config = cfg

	context.SetAccessTokenLock(new(sync.RWMutex))
	context.SetJsAPITicketLock(new(sync.RWMutex))

	//create http client
	if cfg.SslCertFilePath != "" && cfg.SslKeyFilePath != "" {
		if client, err := util.NewTLSHttpClient(cfg.SslCertFilePath, cfg.SslKeyFilePath); err == nil {
			context.SHTTPClient = client
		}
	}

	if cfg.SslCertContent != "" && cfg.SslKeyContent != "" {
		if client, err := util.NewTLSHttpClientFromContent(cfg.SslCertContent, cfg.SslKeyContent); err == nil {
			context.SHTTPClient = client
		}
	}
}

//Mch 商户平台
func (wc *Wechat) Mch() (mch *MchMgr, err error) {
	mch = new(MchMgr)
	mch.Wechat = *wc
	return
}

//Mp 公众平台
func (wc *Wechat) Mp() (mp *MpMgr, err error) {
	err = wc.checkCfgBase()
	if err != nil {
		return
	}
	mp = new(MpMgr)
	mp.Wechat = *wc
	return
}

//checkCfgBase 检查配置基本信息
func (wc *Wechat) checkCfgBase() (err error) {
	if wc.Context.AppID == "" {
		return fmt.Errorf("%s", "配置中没有AppID")
	}
	if wc.Context.AppSecret == "" {
		return fmt.Errorf("%s", "配置中没有AppSecret")
	}
	if wc.Context.Token == "" {
		return fmt.Errorf("%s", "配置中没有Token")
	}
	return
}

func (wc *Wechat) checkCfgMch() (err error) {
	err = wc.checkCfgBase()
	if err != nil {
		return
	}
	if wc.Context.MchID == "" {
		return fmt.Errorf("%s", "配置中没有MchID")
	}
	if wc.Context.MchAPIKey == "" {
		return fmt.Errorf("%s", "配置中没有MchAPIKey")
	}
	if wc.Context.SslCertFilePath == "" || wc.Context.SslCertContent == "" {
		return fmt.Errorf("%s", "配置中没有SslCert")
	}
	if wc.Context.SslKeyFilePath == "" || wc.Context.SslKeyContent == "" {
		return fmt.Errorf("%s", "配置中没有SslKey")
	}
	return
}
