'use strict';

angular
    .module('Quiz.version.directive', [])
    .directive('appVersion', ['version', function(version) {
        return function(scope, ele, attrs) {
            ele.text(version);
        };
    }])
;
