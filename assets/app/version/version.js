'use strict';

angular
    .module('Quiz.version', [
        'Quiz.version.filter',
        'Quiz.version.directive'
    ])

    .value('version', '0.1')
;
