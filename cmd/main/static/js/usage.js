(function () {
    'use strict'

    const dateField = document.getElementById("billDate");

    const now = new Date();
    const tzOffset = (new Date()).getTimezoneOffset() * 60000;
    const lastDay = new Date(now.getFullYear(), now.getMonth() + 1, 0);
    const lastDateLocalISOTime = (new Date(lastDay - tzOffset)).toISOString().substring(0,10);

    dateField.setAttribute('max', lastDateLocalISOTime);

})()