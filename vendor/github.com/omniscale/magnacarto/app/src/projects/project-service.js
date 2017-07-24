angular.module('magna-app')
/* Todo rename to ProjectServicev */
.provider('ProjectService', [function() {
  this.$get = ['$http', '$rootScope', '$q', '$timeout', '$websocket', 'magnaConfig', 'StyleService', 'LayerService', 'DashboardService',
    function($http, $rootScope, $q, $timeout, $websocket, magnaConfig, StyleService, LayerService, DashboardService) {
      var ProjectServiceInstance = function() {
        this.project = undefined;
        this.mmlData = undefined;
        this.dashboardMaps = [];
        this.bookmarkedMaps = [];
        this.socketUrl = undefined;
        this.socket = undefined;
        this.mmlLoadPromise = undefined;
        this.mcpLoadPromise = undefined;
        this.projectLoadedPromise = undefined;
        this.mcpSaveTimeout = undefined;
      };

      ProjectServiceInstance.prototype.loadProject = function(project) {
        var self = this;

        self.unloadProject();

        self.project = project;

        self.mmlLoadPromise = self.loadMML();
        self.mcpLoadPromise = self.loadMCP();

        self.projectLoadedPromise = $q.all([self.mmlLoadPromise, self.mcpLoadPromise]).then(function(data){
          self.handleMMLResponse(data[0].data);
          self.handleMCPResponse(data[1].data);

          self.bindSocket();
          self.enableWatchers();
        });

        return self.projectLoadedPromise;
      };

      ProjectServiceInstance.prototype.loadMML = function() {
        var self = this;
        return $http.get(magnaConfig.projectBaseUrl + self.project.base + '/' + self.project.mml);
      };

      ProjectServiceInstance.prototype.loadMCP = function() {
        var self = this;
        return $http.get(magnaConfig.projectBaseUrl + self.project.base + '/' + self.project.mcp);
      };

      ProjectServiceInstance.prototype.handleMMLResponse = function(response) {
        var self = this;
        if(self.mmlData !== undefined) {
          // Clear array but keep reference to it.
          // If a = [] is used instead of a.length = 0, reference changes
          self.mmlData.Stylesheet.length = 0;
          angular.forEach(response.Stylesheet, function(style) {
            self.mmlData.Stylesheet.push(style);
          });
          self.mmlData.Layer.length = 0;
          angular.forEach(response.Layer, function(layer) {
            self.mmlData.Layer.push(layer);
          });
        } else {
          self.mmlData = response;
        }

        StyleService.setStyles(self.project.available_mss);
        StyleService.setProjectStyles(self.mmlData.Stylesheet);

        LayerService.setLayers(self.mmlData.Layer);
      };

      ProjectServiceInstance.prototype.handleMCPResponse = function(response) {
        var self = this;
        response.dashboardMaps = response.dashboardMaps || [];
        response.bookmarkedMaps = response.bookmarkedMaps || [];
        self.mcpData = response;

        // assign to object property for easy access from outside;
        self.bookmarkedMaps = self.mcpData.bookmarkedMaps;
        DashboardService.maps = self.mcpData.dashboardMaps;
      };

      ProjectServiceInstance.prototype.unloadProject = function() {
        var self = this;
        if(self.mmlData === undefined) {
          return;
        }

        self.disableWatchers();
        if(self.socket !== undefined) {
          self.socket.$close();
        }

        self.project = undefined;
        self.mmlData = undefined;
        self.dashboardMaps = [];
        self.bookmarkedMaps = [];
        self.socketUrl = undefined;
        self.socket = undefined;
        self.mmlLoadPromise = undefined;

        DashboardService.maps = [];
        StyleService.setStyles([]);
        StyleService.setProjectStyles([]);
        LayerService.setLayers([]);
      };

      ProjectServiceInstance.prototype.saveMML = function() {
        var self = this;
        $http.post(magnaConfig.projectBaseUrl + self.project.base + '/' + self.project.mml, angular.toJson(self.mmlData, true));
      };

      ProjectServiceInstance.prototype.saveMCP = function() {
        var self = this;
        if(self.mcpSaveTimeout !== undefined) {
          $timeout.cancel(self.mcpSaveTimeout);
        }
        // prevent too often safe. mostly triggered by gridster when resize or dragging a map
        self.mcpSaveTimeout = $timeout(function() {
          $http.post(magnaConfig.projectBaseUrl + self.project.base + '/' + self.project.mcp, angular.toJson(self.mcpData, true));
          self.mcpSaveTimeout = undefined;
        }, 1000);
      };

      ProjectServiceInstance.prototype.bindSocket = function() {
        var self = this;
        self.socketUrl = angular.copy(magnaConfig.socketUrl);
        self.socketUrl += 'mml=' + self.project.mml;
        self.socketUrl += '&mss=' + self.project.available_mss;
        self.socketUrl += '&base=' + self.project.base;
        self.socket = $websocket.$new({
          url: self.socketUrl,
          reconnect: true,
          reconnectInterval: 100
        });

        self.projectLoadedPromise = self.projectLoadedPromise.then(function() {
          self.socket.$on('$message', function (resp) {
            if(resp.updated_mml === true) {
              self.mmlLoadPromise = self.loadMML();
              self.projectLoadedPromise = $q.all([self.mmlLoadPromise, self.mcpLoadPromise]).then(function(data) {
                self.handleMMLResponse(data[0].data);
              });
            }
          });
        });
      };

      ProjectServiceInstance.prototype.projectLoaded = function() {
        var self = this;
        return self.projectLoadedPromise;
      };

      ProjectServiceInstance.prototype.getSocket = function() {
        return this.socket;
      };

      ProjectServiceInstance.prototype.enableWatchers = function() {
        var self = this;

        // listen on changes in dashboardMaps
        // save if change occurs
        self.dashboardMapsWatcher = $rootScope.$watch(function() {
          return self.mcpData.dashboardMaps;
        }, function(n, o) {
          if(n === o) return;
          self.saveMCP();
        }, true);

        // listen on changes in bookmarkedMaps
        // save if change occurs
        self.bookmarkedMapsWatcher = $rootScope.$watch(function() {
          return self.bookmarkedMaps;
        }, function(n, o) {
          if(n === o) return;
          self.saveMCP();
        }, true);

        self.stylesWatcher = $rootScope.$watch(function() {
          return self.mmlData.Stylesheet;
        }, function(n, o) {
          if(n === o) return;
          self.saveMML();
        }, true);

        self.layersWatcher = $rootScope.$watch(function() {
          return self.mmlData.Layer;
        }, function(n, o) {
          if(n === o) return;
          self.saveMML();
        }, true);
      };

      ProjectServiceInstance.prototype.disableWatchers = function() {
        var self = this;
        if(self.dashboardMapsWatcher !== undefined) {
          self.dashboardMapsWatcher();
          self.dashboardMapsWatcher = undefined;
        }

        if(self.bookmarkedMapsWatcher !== undefined) {
          self.bookmarkedMapsWatcher();
          self.bookmarkedMapsWatcher = undefined;
        }

        if(self.stylesWatcher !== undefined) {
          self.stylesWatcher();
          self.stylesWatcher = undefined;
        }

        if(self.layersWatcher !== undefined) {
          self.layersWatcher();
          self.layersWatcher = undefined;
        }
      };

      return new ProjectServiceInstance();
  }];
}]);
