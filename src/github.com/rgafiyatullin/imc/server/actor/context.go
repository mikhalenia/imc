package actor

import (
	"github.com/rgafiyatullin/imc/server/actor/logging"
	"github.com/rgafiyatullin/imc/server/actor/metrics"
	"github.com/rgafiyatullin/imc/server/config"
	"os"
)

type Ctx interface {
	Path() string
	Log() logging.Ctx
	Metrics() metrics.Ctx
	NewChild(name string) Ctx
	Halt(code int, message string)
}

type impl struct {
	path_    string
	log_     logging.Ctx
	metrics_ metrics.Ctx
}

func (this *impl) Log() logging.Ctx     { return this.log_ }
func (this *impl) Metrics() metrics.Ctx { return this.metrics_ }
func (this *impl) Path() string         { return this.path_ }

func (this *impl) Halt(code int, message string) {
	this.log_.Fatal("System halt requested. [code: %d; msg: %s]", code, message)
	this.log_.Flush()
	os.Exit(code)
}

func (this *impl) NewChild(name string) Ctx {
	child := new(impl)
	child.metrics_ = this.metrics_
	childName := this.path_ + "/" + name
	child.path_ = childName
	child.log_ = this.log_.CloneWithName(childName)
	return child
}

func New(config config.Config) Ctx {
	ctx := new(impl)
	ctx.log_ = logging.New(config)
	ctx.metrics_ = metrics.New(ctx.log_.CloneWithName("/metrics"), config)
	ctx.path_ = ""
	return ctx
}
