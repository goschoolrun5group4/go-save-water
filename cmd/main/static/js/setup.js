(function () {
    'use strict'

    const dateField = document.getElementById("billDate");

    const now = new Date();
    const tzOffset = (new Date()).getTimezoneOffset() * 60000;

    const firstDay = new Date(now.getFullYear(), now.getMonth(), 1);
    const firstDateLocalISOTime = (new Date(firstDay - tzOffset)).toISOString().substring(0,10);

    const lastDay = new Date(now.getFullYear(), now.getMonth() + 1, 0);
    const lastDateLocalISOTime = (new Date(lastDay - tzOffset)).toISOString().substring(0,10);

    dateField.setAttribute('min', firstDateLocalISOTime)
    dateField.setAttribute('max', lastDateLocalISOTime)

})()