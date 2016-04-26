'use strict';

angular
    .module('Quiz')

    .controller('QuizCtrl.Login', ['$scope', '$cookies', '$http', '$location', function($scope, $cookies, $http, $location) {
        console.log('QuizCtrl.Login:', $scope);
        $scope.logging = false;
        $scope.login = function() {
            //console.log('QuizCtrl.Login: ', this);
            let enc = encodeURIComponent;
            let name = enc(this.username);
            let pass = enc(this.password);
            var that = this;
            console.log('QuizCtrl.Login: ', name, pass);
            this.logging = true;
            this.response = null;
            $http({
                url: '/user', method: "POST",
                headers: {'Content-Type': 'application/x-www-form-urlencoded'},
                data:'method=login&email='+name+'&pass='+pass,
                dataType: 'json'
            }).success(function(response){
                console.log('QuizCtrl.Login: ', response);
                that.response = response
                if (response.code == 0) {
                    $location.path('/exam');
                    $cookies.put("token", response.token); // TODO: good auth
                    that.password = '';
                }
                that.logging = false;
            });
        };
    }])

    .config(['$routeProvider', function($route) {
        $route 
            .when('/login', {
                templateUrl: '/v/login',
                controller: 'QuizCtrl.Login'
            })
        ;
    }])
;
