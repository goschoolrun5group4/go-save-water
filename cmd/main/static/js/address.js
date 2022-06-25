(function () {
    'use strict'

    const btn =  document.getElementById("apiCall");
    let blockNum = document.getElementById("blockNumber");
    let street = document.getElementById("street");
    let buildingName = document.getElementById("buildingName");
    let postalCode = document.getElementById("postalCode");
    let formSubmit = document.getElementById("formSubmit");

    btn.addEventListener('click', (e) => {
        e.preventDefault();

        fetch('https://developers.onemap.sg/commonapi/search?searchVal='+postalCode.value+'&returnGeom=N&getAddrDetails=Y')
            .then(response => {
                return response.json()
            })
            .then(json => {
                return json['results'][0];
            })
            .then(result => {
                if (result != undefined) {
                    postalCode.classList.remove("is-invalid");
                    formSubmit.disabled = false;
                    blockNum.value = result["BLK_NO"];
                    street.value = result["ROAD_NAME"];
                    if (result["BUILDING"] != 'NIL') {
                        buildingName.value = result["BUILDING"];
                    } else {
                        buildingName.value = "";
                    };
                } else {
                    postalCode.classList.add('is-invalid');
                    formSubmit.disabled = true;
                };
            })
    });

})()