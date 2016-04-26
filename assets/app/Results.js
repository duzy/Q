'use strict';

angular
    .module('Quiz')

    .controller('QuizCtrl.Results', [function() {
        console.log('QuizCtrl.Results:', this);
    }])

    .config(['$routeProvider', function($route) {
        $route
            .when('/results', {
                templateUrl: '/v/results',
                controller: 'QuizCtrl.Results'
            })
        ;
    }])
;
