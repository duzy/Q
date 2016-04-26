'use strict';

angular
    .module('Quiz')

    .controller('QuizCtrl.Marks', [function() {
        console.log('QuizCtrl.Marks:', this);
    }])

    .config(['$routeProvider', function($route) {
        $route
            .when('/marks', {
                templateUrl: '/v/marks',
                controller: 'QuizCtrl.Marks'
            })
        ;
    }])
;
