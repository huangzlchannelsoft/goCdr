// mon
package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/client_golang/prometheus/push"
)

const (
	METRIX_GAUGE    = 0
	METRIX_GAUGEVEC = 1
)

type NewMetrix func(job string, typ int, prex string, name string, labels []string)
type SetMetrix func(val float64, labelsval []string, job string, idx int)
type RegMetrix func(job string)

type InitPluginFunc func(args ...string)
type SetPluginMetaFunc func(newMetrix NewMetrix, regMetrix RegMetrix)
type CallPluginFunc func(setMetrix SetMetrix)

type Plugin struct {
	name       string
	Init       InitPluginFunc
	GetMeta    SetPluginMetaFunc
	Call       CallPluginFunc
	metrixType []int
	pusher     *push.Pusher
	metrix     []prometheus.Collector
}

var gPlugins []*Plugin

func init() {
	log.Println("mon!")
	gPlugins = append(gPlugins, &Plugin{
		"cdrCallStat",
		InitCallStat,
		SetCallStatMeta,
		CallCallStat,
		nil, //metrixType  []int
		nil, //pusher      *push.Pusher
		nil, //metrix      []prometheus.Collector
	})
}

func PromethuesClient(pushOrPull bool, pushOrPullUri string, nid string, sampleSec int) {
	log.Println("enter promethuesClient.")

	isPull := func() bool {
		return !pushOrPull
	}

	var thisPlugin *Plugin
	regMetrix := func(job string) {
		if thisPlugin != nil {
			if isPull() {
				prometheus.MustRegister(thisPlugin.metrix...)
			} else {
				registry := prometheus.NewRegistry()
				registry.MustRegister(thisPlugin.metrix...)
				thisPlugin.pusher = push.New(pushOrPullUri, job).Gatherer(registry)
				if nid != "" {
					thisPlugin.pusher = thisPlugin.pusher.Grouping("instance", nid)
				}
			}
		}
	}
	newMetrix := func(job string, typ int, prex string, name string, labels []string) {
		if thisPlugin != nil {
			thisPlugin.metrixType = append(thisPlugin.metrixType, typ)
			switch typ {
			case METRIX_GAUGE:
				thisPlugin.metrix = append(thisPlugin.metrix, prometheus.NewGauge(
					prometheus.GaugeOpts{
						Name: fmt.Sprintf("%s_%s", prex, name),
						Help: "",
					}))
			case METRIX_GAUGEVEC:
				exlabels := labels
				if isPull() {
					exlabels = append(exlabels, []string{"job", "instance"}...)
				}
				thisPlugin.metrix = append(thisPlugin.metrix, prometheus.NewGaugeVec(
					prometheus.GaugeOpts{
						Name: fmt.Sprintf("%s_%s", prex, name),
						Help: "",
					},
					exlabels,
				))
			}
		}
	}
	setMetrix := func(val float64, labelsval []string, job string, idx int) {
		if thisPlugin != nil {
			switch thisPlugin.metrixType[idx] {
			case METRIX_GAUGE:
				thisPlugin.metrix[idx].(prometheus.Gauge).Set(val)
			case METRIX_GAUGEVEC:
				exlabelsval := labelsval
				if isPull() {
					exlabelsval = append(exlabelsval, []string{job, nid}...)
				}
				thisPlugin.metrix[idx].(*prometheus.GaugeVec).WithLabelValues(exlabelsval...).Set(val)
			}
		}
	}
	// for _, pcfg := range gCfg.Plugins {
	// 	found := false
	// 	for _, plugin := range gPlugins {
	// 		if pcfg.Name == plugin.name {
	// 			found = true
	// 			plugin.Init(pcfg.Params...)
	// 			thisPlugin = plugin
	// 			plugin.GetMeta(newMetrix, regMetrix)
	// 		}
	// 	}
	// 	if !found {
	// 		log.Println("[Err] not found plugin:", pcfg.Name)
	// 	}
	// }

	for _, plugin := range gPlugins {
		plugin.Init()
		thisPlugin = plugin
		plugin.GetMeta(newMetrix, regMetrix)
	}

	if isPull() {
		go func() {
			http.Handle("/metrics", promhttp.Handler())
			log.Fatal(http.ListenAndServe(pushOrPullUri, nil))
		}()
	}

	sampleTicker := time.NewTicker(time.Duration(sampleSec) * time.Second)
	for {
		select {
		case <-sampleTicker.C:
			for _, plugin := range gPlugins {
				if plugin.metrix == nil {
					continue
				}
				thisPlugin = plugin
				plugin.Call(setMetrix)
				if plugin.pusher != nil {
					err := plugin.pusher.Add()
					if err != nil {
						log.Println("[Err] mon push.", err.Error())
					}
				}
			}
		}
	}
}
