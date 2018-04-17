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
    rpc('flipflop.FlipFlop', {}, function (data) {
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

var ticTacToeGameID;
var ticTacToePlayerID;
var ticTacToeSetintervalID;

function ticTacToeJoin() {
    rpc('tictactoe.Join', {}, function (data) {
        console.log('Success');
        console.log(data);
        ticTacToeGameID = data.result.GameID;
        ticTacToePlayerID = data.result.PlayerID;
        ticTacToeRedraw(data.result.Board);
        ticTacToeCheckWinner(data.result);
        if (ticTacToeSetintervalID == null) {
            ticTacToeSetintervalID = setInterval(ticTacToeGetGame, 3000);
        }
    }, alertFail);
}

function ticTacToePlace(row, col) {
    rpc('tictactoe.Place', {
        GameID: ticTacToeGameID,
        PlayerID: ticTacToePlayerID,
        Position: [row, col],
    }, function (data) {
        console.log('Success');
        console.log(data);
        ticTacToeRedraw(data.result.Board);
        ticTacToeCheckWinner(data.result);
    }, alertFail);
}

function ticTacToeGetGame() {
    rpc('tictactoe.GetGame', { GameID: ticTacToeGameID }, function (data) {
        console.log('Success');
        console.log(data);
        ticTacToeRedraw(data.result.Board);
        ticTacToeCheckWinner(data.result);
    }, function (data) {
        clearInterval(ticTacToeSetintervalID);
    });
}

var ticTacToePieces = {
    '0': '?',
    '1': 'X',
    '-1': 'O',
};
function ticTacToeRedraw(board) {
    board.forEach(function (row, i) {
        row.forEach(function (val, j) {
            $(`#row${i}col${j}`)[0].textContent = ticTacToePieces[val];
        });
    });
}

function ticTacToeCheckWinner(result) {
    if (!result.Over) {
        return;
    }
    clearInterval(ticTacToeSetintervalID);
    switch (result.Winner) {
        case '':
            alert('Tic-tac-toe is a draw!');
            break;
        case ticTacToePlayerID:
            alert('You won tic-tac-toe!');
            break;
        default:
            alert('You lost tic-tac-toe!');
    }
}

$(document).ready(function () { });
