'use strict';

angular
    .module('Quiz')

    .controller('QuizCtrl.Exam', ['$scope', '$routeParams', '$cookies', '$http', '$location', function($scope, $params, $cookies, $http, $location) {
        //console.log('QuizCtrl.Exam:', $params, $scope);

        $scope.number = Number.parseInt($params.n || '1');
        $scope.submitting = false;
        $scope.submit = function() {
            //console.log('QuizCtrl.Exam:', this);
            let enc = encodeURIComponent;
            let t = $cookies.get('token');
            let a = [], i = 0, s = '';
            if (this.answer) {
                a.push(enc(this.answer));
            } else {
                for (i=1; i <= 10; ++i) {
                    s = 'answer' + i;
                    if (this[s]) a.push(this[s]); else break;
                }
            }
            
            if (a.length == 0) {
                this.response = { code:-1, message: "You need to pick an answer!" };
                return;
            }
            
            s = JSON.stringify({a:a});
            i = this.number;

            console.log('QuizCtrl.Exam: answer:', s, a);

            this.submitting = true;
            var that = this;
            $http({
                url: '/take', method: "POST",
                headers: {'Content-Type': 'application/x-www-form-urlencoded'},
                data:'t='+t+'&n='+i+'&a='+s, dataType: 'json'
            }).success(function(response){
                console.log('QuizCtrl.Exam: ', response);
                that.response = response;
                if (response.code === 0) {
                    $location.url('/exam?n='+(that.number += 1));
                }
                that.submitting = false;
            });
        };
        
        function init() {
            //var m = $params.m;
            //$location.url('/exam?m=' + (m||""))
            //$location.path('/exam/' + (m||""))
        }
        
        init()
    }])

    .config(['$routeProvider', function($route) {
        $route
            .when('/exam', { // Note: using RESTful and JSON might be better.
                controller: 'QuizCtrl.Exam',
                templateUrl: function(p) {
                    /*
                    let t = $cookies.get('token');
                    console.log('QuizCtrl.Exam: /exam: ', p, t)
                    return t ? ('/v/exam?n=' + (p.n || 1)) : '/';
                    */
                    return '/v/exam?n=' + (p.n || 1);
                }
            })
        ;
    }])
;
