'use strict';

angular
    .module('Quiz', [
        'ngRoute', 
        'ngCookies',
        //'Quiz.Login',
        //'Quiz.Logout',
        //'Quiz.Register',
        //'Quiz.Exam',
        //'Quiz.Results',
        //'Quiz.Marks',
        'Quiz.version'
    ])

    .controller('QuizCtrl', [function() {
        console.log("QuizCtrl:", this);
    }])

    .config(['$routeProvider', '$locationProvider', function($route, $location) {
        $route
            .when('/', {
                //templateUrl: '/hi',
                //controller: 'QuizCtrl'
                redirectTo: '/login'
            })
            .otherwise({ redirectTo: '/login' })
        ;
    }])

    .run(['$rootScope', '$cookieStore', function($scope, $cookie) {
        console.log("Quiz.run: $scope:",     $scope);
        console.log("Quiz.run: $cookie:",    $cookie);
    }])
;
