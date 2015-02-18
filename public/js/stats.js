'use strict'

var StatsController = {
    chartContainer: $('#container'),

    init: function() {
        this.series = this.parseHash(document.TempData);
        this.period = this.getUrlParam("period");
        if (this.period == "") {
            this.period = "5m"
        }
        $('.period').each(function(i, el){
            if($(el).html() == StatsController.period) {
                $(el).addClass('active');
            }
        })
        this.initHighCharts();
    },

    parseHash: function(hash) {
        var data = [];
        for (var k in hash) {
            var date = Date.parse(k);
            data.push([date, parseFloat(hash[k])]);
        }
        return data;
    },

    getData: function(lastDate) {
        $.get('/arduino/last_stats', {lastDate: lastDate.toISOString(), period: StatsController.period}, function(data){
            StatsController.stats = StatsController.parseHash(data);
        });
    },

    loadSeries: function() {
        var that = this;
        setInterval(function () {
            var series = that.series[0];
            var lastDate;
            if (series.data.length > 0) {
                lastDate = new Date(series.data[series.data.length-1].x);
            } else {
                lastDate = new Date("1970-01-01");
            }
            StatsController.getData(lastDate);
            var data = StatsController.stats;
            if (!data || data.length < 1) { return };
            for (var i = 0; i < data.length; i++) {
                var value = data[i];
                series.addPoint([value[0], value[1]], true, true);
            }
        }, 1000);
    },

    initHighCharts: function() {
        this.chart = this.chartContainer.highcharts({
            chart: {
                type: 'spline',
                animation: Highcharts.svg,
                events: { load: this.loadSeries }
            },
            title: { text: 'Temperature Stats' },
            subtitle: { text: 'Autohome' },
            xAxis: { type: 'datetime' },
            yAxis: { title: { text: 'Temperature (°C)' }, min: 20 },
            plotOptions: { line: { enableMouseTracking: true } },
            series: [{ name: 'Temp', data: this.series }]
        });
    },

    getUrlParam: function(name) {
        name = name.replace(/[\[]/, "\\[").replace(/[\]]/, "\\]");
        var regex = new RegExp("[\\?&]" + name + "=([^&#]*)"),
            results = regex.exec(location.search);
        return results === null ? "" : decodeURIComponent(results[1].replace(/\+/g, " "));
    },
};

$(function () {
    StatsController.init()
});
