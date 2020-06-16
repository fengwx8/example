package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"time"
)

var (
	requestCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name:      "request_total",
			Help:      "Number of request processed by this service.",
		}, []string{},
	)

	requestLatency = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:      "request_latency_seconds",
			Help:      "Time spent in this service.",
			Buckets:   []float64{0.01, 0.02, 0.05, 0.1, 0.2, 0.5, 1.0, 2.0, 5.0, 10.0, 20.0, 30.0, 60.0, 120.0, 300.0},
		}, []string{},
	)
	//带宽可用率，13：00最低，百分比
	bitrate := prometheus.NewGauge(prometheus.GaugeOpts{
		Name:      "bitrate",
		Help:      "Bitrate availability for service.",
	})
	//服务器繁忙程度，百分比，以可用率为度量
	getstate := prometheus.NewGauge(prometheus.GaugeOpts{
		Name:      "server_state",
		Help:      "Busy rate of server.",
	})
)

// AdmissionLatency measures latency / execution time of Admission Control execution
// usual usage pattern is: timer := NewAdmissionLatency() ; compute ; timer.Observe()
type RequestLatency struct {
	histo *prometheus.HistogramVec
	start time.Time
}

func Register() {
	prometheus.MustRegister(requestCount)
	prometheus.MustRegister(requestLatency)
	prometheus.MustRegister(bitrate)
	prometheus.MustRegister(getstate)
}


// NewAdmissionLatency provides a timer for admission latency; call Observe() on it to measure
func NewAdmissionLatency() *RequestLatency {
	return &RequestLatency{
		histo: requestLatency,
		start: time.Now(),
	}
}

// Observe measures the execution time from when the AdmissionLatency was created
func (t *RequestLatency) Observe() {
	(*t.histo).WithLabelValues().Observe(time.Now().Sub(t.start).Seconds())
}


// RequestIncrease increases the counter of request handled by this service
func RequestIncrease() {
	requestCount.WithLabelValues().Add(1)
	hour:=time.Now().Hour()
	minute:=time.Now().Minute()
	//转化成10进制对应的值，相当于小数部分
	minutetopoint:=float64(minute)/60
	hours:=float64(hour)
	//根据时间模拟流量，这也比较符合事实
	//简单地用顶点为(13,10)的抛物线模拟比特率随时间变化的对应关系，则最低值为(13,10)，最高值为(0,94.5)
	//float64 BdW = (hours + minutetopoint - 13) * (hours + minutetopoint - 13) / 2 + 10
	BdW:=float64(10)
	BdW=BdW+(hours + minutetopoint - 13) * (hours + minutetopoint - 13) / 2
	bitrate.Set(BdW)
	//直接将可用率度量服务器的繁忙程度
	getstate.Set(100 - BdW)

}
