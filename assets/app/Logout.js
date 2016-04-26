'use strict';

angular
    .module('Quiz')

    .controller('QuizCtrl.Logout', ['$scope', '$http', function() {
        console.log('QuizCtrl.Logout:', $scope);
    }])

    .config(['$routeProvider', function($route) {
        $route
            .when('/logout', {
                templateUrl: '/v/logout',
                controller: 'QuizCtrl.Logout'
            })
        ;
    }])
;
