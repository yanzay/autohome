$(function () {

    parseHash = function(hash) {
        var data = [];
        for (var k in hash) {
            date = Date.parse(k);
            data.push([date, parseFloat(hash[k])]);
        }
        return data
    }

    series = parseHash(document.TempData)

    $('#container').highcharts({
        chart: {
            type: 'spline',
            animation: Highcharts.svg,
            events: {
                load: function () {
                    that = this
                    setInterval(function () {
                        var series = that.series[0];
                        var lastDate = new Date(series.data[series.data.length-1].x);
                        $.get('/arduino/last_stats', {lastDate: lastDate.toISOString()}, function(data){
                            var d = parseHash(data);
                            for (var i = 0; i < d.length; i++) {
                                value = d[i];
                                series.addPoint([value[0], value[1]], true, true);
                            }
                        })
                    }, 1000);
                }
            }
        },
        title: {
            text: 'Temperature Stats'
        },
        subtitle: {
            text: 'Autohome'
        },
        xAxis: {
            type: 'datetime'
        },
        yAxis: {
            min: 20,
            title: {
                text: 'Temperature (°C)'
            }
        },
        plotOptions: {
            line: {
                enableMouseTracking: true
            }
        },
        series: [{
            name: 'Temp',
            data: series
        }]
    });
});
