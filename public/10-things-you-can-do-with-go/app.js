function rpc(method, params, done, fail) {
    $.ajax({
        type: 'POST',
        url: '/rpc/dowithgo',
        contentType: 'application/json',
        data: JSON.stringify({
            jsonrpc: '2.0',
            id: 0,
            method: method,
            params: [params],
        }),
    }).done(done).fail(fail);
}

function alertFail(data) {
    alert('Error: ' + data.responseJSON.error);
}

function flipFlop() {
    rpc('flipflop.FlipFlopr', {}, function (data) {
        console.log('Success');
        console.log('Flipped: ' + data.result.Flipped);
        console.log('Count: ' + data.result.Count);
    }, alertFail);
}

function shuffle() {
    rpc('shuffle.Shuffle', {}, function (data) {
        console.log('Success');
        console.log(data);
    }, alertFail);
}

$(document).ready(function () { });
