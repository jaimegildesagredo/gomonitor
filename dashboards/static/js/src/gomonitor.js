(function () {
  // TODO: Separate repository and service
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

      that.ToggleHistory = function (value) {
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

  var BandwidthChart = function (options) {
    var that = {},
        parentElement = d3.select(options.parentSelector),
        margin = options.margin,
        width = options.width - margin.left - margin.right,
        height = options.height - margin.top - margin.bottom,
        svg = parentElement.append("svg")
                             .attr("width", width + margin.left + margin.right)
                             .attr("height", height + margin.top + margin.bottom)
                           .append("g")
                             .attr("transform", "translate(" + margin.left + ", " + margin.top + ")"),
        x = d3.time.scale().range([0, width]),
        y = d3.scale.linear().range([height, 0]),
        xAxis = d3.svg.axis().scale(x).orient("bottom").ticks(5),
        yAxis = d3.svg.axis().scale(y).orient("left").ticks(5),
        upLine = d3.svg.line(),
        downLine = d3.svg.line();

    upLine.x(function (bandwidth) { return x(bandwidth.created_at); });
    upLine.y(function (bandwidth) { return y(bandwidth.up); });

    downLine.x(function (bandwidth) { return x(bandwidth.created_at); });
    downLine.y(function (bandwidth) { return y(bandwidth.down); });

    that.draw = function (bandwidths) {
      var upAndDownBandwidths = [];
      upAndDownBandwidths = upAndDownBandwidths.concat(bandwidths.map(function (bandwidth) { return bandwidth.up; }));
      upAndDownBandwidths = upAndDownBandwidths.concat(bandwidths.map(function (bandwidth) { return bandwidth.down; }));

      x.domain(d3.extent(bandwidths, function (bandwidth) { return bandwidth.created_at; }));
      y.domain([0, d3.max(upAndDownBandwidths, function (bandwidth) { return bandwidth; })]);

      svg.selectAll("path").remove();
      svg.append("path")
           .attr("class", "upLine")
           .attr("d", upLine(bandwidths));

      svg.append("path")
           .attr("class", "downLine")
           .attr("d", downLine(bandwidths));

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

  var main = function () {
    var monitorBandwidth = function (interfaceName) {
      var networkBandwidthList = document.getElementById("network-bandwidth");
      var itemElement = document.createElement("li"),
          containerElement = document.createElement("div"),
          parentElement = document.createElement("div"),
          titleElement = document.createElement("h3");

      titleElement.appendChild(document.createTextNode(interfaceName));
      containerElement.appendChild(titleElement);
      parentElement.setAttribute("id", interfaceName + "-graph");
      containerElement.appendChild(parentElement);
      itemElement.appendChild(containerElement);
      networkBandwidthList.appendChild(itemElement);

      var bandwidthChart = BandwidthChart({
        parentSelector: "#" + interfaceName + "-graph",
        width: 640,
        height: 480,
        margin: {
          top: 30,
          right: 20,
          bottom: 30,
          left: 75
        }
      });

      interfacesRepository.onBandwidth(interfaceName, function () {
        bandwidthChart.draw(interfacesRepository.findAllBandwidths(interfaceName));
      })
      interfacesRepository.monitorBandwidth(interfaceName);
    }

    var interfacesRepository = InterfacesRepository({
      baseUrl: "http://" + document.location.hostname  + ":3000/networks",
      historyEnabled: false,
      historyLimit: 25
    });

    interfacesRepository.findAll(function (interfaces) {
      for (var i=0; i<interfaces.length; i++) {
        if (interfaces[i].state != "down") {
          monitorBandwidth(interfaces[i].name);
        }
      }
    });

    var enableHistoryLimitButton = document.getElementById('enable-history-button');
    enableHistoryLimitButton.addEventListener("click", function (e) {
      e.preventDefault();
      interfacesRepository.ToggleHistory();
      e.target.classList.toggle("activated");
    }, false);
  }

  main();
})()
