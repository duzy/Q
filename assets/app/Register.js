'use strict';

angular
    .module('Quiz')

    .controller('QuizCtrl.Register', ['$scope', '$cookies', '$http', '$location', function($scope, $cookies, $http, $location) {
        console.log('QuizCtrl.Register:', $scope);
        $scope.regging = false
        $scope.register = function() {
            let enc = encodeURIComponent;
            let name = enc(this.user.username);
            let pass = enc(this.user.password);
            let first = enc(this.user.firstName);
            let last = enc(this.user.lastName);
            var that = this;
            console.log('QuizCtrl.Register: ', this.user);
            this.regging = true
            this.response = null;
            $http({
                url: '/user', method: "POST",
                headers: {'Content-Type': 'application/x-www-form-urlencoded'},
                data:'method=register&email='+name+'&pass='+pass+'&firstName='+first+'&lastName='+last,
                dataType: 'json'
            }).success(function(response){
                console.log('QuizCtrl.Register: ', response);
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
            .when('/register', {
                templateUrl: '/v/register',
                controller: 'QuizCtrl.Register'
            })
        ;
    }])
;
