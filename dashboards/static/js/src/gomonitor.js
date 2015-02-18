(function () {
  var COLORS = [
    "red",
    "blue",
    "green"
  ]

  var LoadRepository = function (options) {
    var that = {},
        baseUrl = options.baseUrl,
        historyEnabled = options.historyEnabled,
        historyLimit = options.historyLimit,
        onLoadCallbacks = [],
        data = [];

    that.findAll = function () {
      return data;
    }

      that.toggleHistory = function (value) {
        historyEnabled = !historyEnabled;
      }

    that.onLoad = function (callback) {
      onLoadCallbacks.push(callback);
    }

    that.monitor = function () {
      window.setInterval(function () { sendRequest(); }, 1000);
    }

    var sendRequest = function (interfaceName) {
      var request = new XMLHttpRequest();
      request.onreadystatechange = function () {
        if (request.readyState === 4) {
          data.push(deserialize(request.response));

          if (!historyEnabled) {
            while (data.length > historyLimit) {
              data.shift();
            }
          }

          for (var o=0; o<onLoadCallbacks.length; o++) {
            onLoadCallbacks[o]();
          }
        }
      }
      request.open("GET", baseUrl);
      request.send(null);
    }

    var deserialize = function (raw) {
      var data = JSON.parse(raw);
      return {
        "created_at": new Date(data.created_at),
        "values": data.values,
      }
    }

    return that;

  }
  var InterfacesRepository = function (options) {
      var that = {},
          baseUrl = options.baseUrl,
          historyEnabled = options.historyEnabled,
          historyLimit = options.historyLimit,
          onBandwidthCallbacks = {},
          data = {};

      that.onBandwidth = function (interfaceName, callback) {
        if (!onBandwidthCallbacks[interfaceName]) {
          onBandwidthCallbacks[interfaceName] = [];
        }
        onBandwidthCallbacks[interfaceName].push(callback);
      }

      that.findAllBandwidths = function (interfaceName) {
        return data[interfaceName] || [];
      }

      that.findAll = function (callback) {
        var request = new XMLHttpRequest();
        request.onreadystatechange = function () {
          if (request.readyState === 4) {
            callback(deserializeInterfaces(request.response));
          }
        }
        request.open("GET", baseUrl);
        request.send(null);
      }

      var deserializeInterfaces = function (raw) {
        return JSON.parse(raw);
      }

      that.toggleHistory = function (value) {
        historyEnabled = !historyEnabled;
      }

      that.monitorBandwidth = function (interfaceName) {
        window.setInterval(function () { sendRequest(interfaceName); }, 1000);
      }

      var sendRequest = function (interfaceName) {
        var request = new XMLHttpRequest();
        request.onreadystatechange = function () {
          if (request.readyState === 4) {
            var interfaceBandwidths;

            if (!data[interfaceName]) {
              data[interfaceName] = [];
            }

            interfaceBandwidths = data[interfaceName];
            interfaceBandwidths.push(deserialize(request.response));

            if (!historyEnabled) {
              while (interfaceBandwidths.length > historyLimit) {
                interfaceBandwidths.shift();
              }
            }

            for (var o=0; o<onBandwidthCallbacks[interfaceName].length; o++) {
              onBandwidthCallbacks[interfaceName][o]();
            }
          }
        }
        request.open("GET", baseUrl + "/" + interfaceName + "/bandwidth");
        request.send(null);
      }

      var deserialize = function (raw) {
        var value = JSON.parse(raw);

        return {
          created_at: new Date(value.created_at),
          down: toKBs(value.down),
          up: toKBs(value.up)
        }
      }

      var toKBs = function (value) {
        return value == 0 ? value : value / 1024;
      }

      return that;
  }

  var LineChart = function (options) {
    var that = {},
        element = d3.select(options.element),
        margin = options.margin,
        width = options.width - margin.left - margin.right,
        height = options.height - margin.top - margin.bottom,
        svg = element.append("svg")
                             .attr("width", width + margin.left + margin.right)
                             .attr("height", height + margin.top + margin.bottom)
                           .append("g")
                             .attr("transform", "translate(" + margin.left + ", " + margin.top + ")"),
        x = d3.time.scale().range([0, width]),
        y = d3.scale.linear().range([height, 0]),
        xAxis = d3.svg.axis().scale(x).orient("bottom").ticks(5),
        yAxis = d3.svg.axis().scale(y).orient("left").ticks(5),
        lines = [];

    var newLine = function () {
      var line = d3.svg.line();
      line.x(function (data) { return x(data[0]); });
      line.y(function (data) { return y(data[1]); });
      return line
    }

    that.draw = function (data) {
      if (lines.length == 0) {
        for (var i=0; i<data.length; i++) {
          lines.push(newLine());
        }
      }

      var flattenedData = [];
      for (var i=0; i<data.length; i++) {
        flattenedData = flattenedData.concat(data[i]);
      }

      x.domain(d3.extent(flattenedData, function (data) { return data[0]; }));
      y.domain([0, d3.max(flattenedData, function (data) { return data[1]; })]);

      svg.selectAll("path").remove();
      for (var i=0; i<data.length; i++) {
        svg.append("path")
             .attr("class", "line " + COLORS[i])
             .attr("d", lines[i](data[i]));
      }

      svg.selectAll("g").remove();
      svg.append("g")
            .attr("class", "x axis")
            .attr("transform", "translate(0, " + height + ")")
            .call(xAxis);
      svg.append("g")
            .attr("class", "y axis")
            .call(yAxis);
    }

    return that;
  }

  var LoadChartFactory = function (element, loadRepository) {
    return LoadChartPresenter({
      view: LoadChartView({
        element: element,
      }),
      loadRepository: loadRepository,
    });
  }

  var LoadChartPresenter = function (options) {
      var that = {},
          loadRepository = options.loadRepository,
          view = options.view;

      loadRepository.onLoad(function () {
        view.render(loadRepository.findAll());
      });

      loadRepository.monitor();

      return that;
  }

  var LoadChartView = function (options) {
    var that = {},
        element = options.element,
        chart = LineChart({
          element: element,
          width: 640,
          height: 480,
          margin: {
            top: 30,
            right: 20,
            bottom: 30,
            left: 75
          }
        });

    that.render = function (data) {
      chart.draw(convertData(data));
    }

    var convertData = function (data) {
      var loadOne = [],
          loadFive = [],
          loadFifteen = [];

      for (var j=0; j<data.length; j++) {
        var item = data[j],
            created_at = item.created_at;

        loadOne.push([created_at, item.values[0]]);
        loadFive.push([created_at, item.values[1]]);
        loadFifteen.push([created_at, item.values[2]]);
      }

      return [loadOne, loadFive, loadFifteen];
    }

    return that;
  }

  var BandwidthChartFactory = function (interfaceName, element, interfacesRepository) {
    return BandwidthChartPresenter({
      interfaceName: interfaceName,
      view: BandwidthChartView({
        element: element,
        interfaceName: interfaceName
      }),
      interfacesRepository: interfacesRepository
    });
  }

  var BandwidthChartPresenter = function (options) {
    var that = {},
        interfaceName = options.interfaceName,
        view = options.view,
        interfacesRepository = options.interfacesRepository;

    interfacesRepository.onBandwidth(interfaceName, function () {
      view.render(interfacesRepository.findAllBandwidths(interfaceName));
    })

    interfacesRepository.monitorBandwidth(interfaceName);

    return that;
  }

  var BandwidthChartView = function (options) {
    var that = {},
        interfaceName = options.interfaceName,
        element = options.element,
        containerElement = document.createElement("div"),
        titleElement = document.createElement("h3"),
        chartElement = document.createElement("div");

    titleElement.appendChild(document.createTextNode(interfaceName));
    containerElement.appendChild(titleElement);
    chartElement.setAttribute("id", interfaceName + "-chart");
    containerElement.appendChild(chartElement);
    element.appendChild(containerElement);

    var chart = LineChart({
      element: chartElement,
      width: 640,
      height: 480,
      margin: {
        top: 30,
        right: 20,
        bottom: 30,
        left: 75
      }
    });

    that.render = function (data) {
      chart.draw(convertData(data));
    }

    return that;
  }

  var convertData = function (data) {
    var up = [],
        down = [];

    for (var i=0; i<data.length; i++) {
      var item = data[i],
          created_at = item.created_at;

      up.push([created_at, item.up]);
      down.push([created_at, item.down]);
    }

    return [up, down];
  }

  var EnableHistoryButtonFactory = function (element, repository) {
    return EnableHistoryButtonPresenter({
      view: EnableHistoryButtonView({element: element}),
      repository: repository,
    });
  }

  var EnableHistoryButtonPresenter = function (options) {
    var that = {},
        enabled = false,
        view = options.view,
        repository = options.repository;

    view.onClick(function () {
      enabled = !enabled;
      if (enabled) {
        view.enabled();
      } else {
        view.disabled();
      }
      repository.toggleHistory();
    });

    return that;
  }

  var EnableHistoryButtonView = function (options) {
    var that = {},
      element = options.element;

    that.onClick = function (callback) {
      element.addEventListener("click", function (e) {
        e.preventDefault();
        callback();
      }, false);
    }

    that.enabled = function () {
      element.textContent = "Disable history"
    }

    that.disabled = function () {
      element.textContent = "Enable history"
    }

    return that;
  }

  var networkDashboard = function () {
    var interfacesRepository = InterfacesRepository({
      baseUrl: "http://" + document.location.hostname  + ":3000/networks",
      historyEnabled: false,
      historyLimit: 25
    });

    var bandwidthCharts = document.getElementById("network-bandwidth");
    interfacesRepository.findAll(function (interfaces) {
      for (var i=0; i<interfaces.length; i++) {
        if (interfaces[i].state != "down") {
          var itemElement = document.createElement("li");
          BandwidthChartFactory(interfaces[i].name, itemElement, interfacesRepository);
          bandwidthCharts.appendChild(itemElement);
        }
      }
    });

    EnableHistoryButtonFactory(document.getElementById("enable-history-button"), interfacesRepository);
  }

  var loadDashboard = function () {
    var loadRepository = LoadRepository({
      baseUrl: "http://" + document.location.hostname  + ":3000/load",
      historyEnabled: false,
      historyLimit: 25
    });

    LoadChartFactory(document.getElementById("load-chart"), loadRepository);
    EnableHistoryButtonFactory(document.getElementById("enable-history-button"), loadRepository);
  }

  var main = function () {
    var dashboardName = window.location.pathname.split("/")[2];

    if (dashboardName == "network") {
      networkDashboard();
    } else if (dashboardName == "load") {
      loadDashboard();
    }
  }
  main();
})()
