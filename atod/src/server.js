const port = 5000

const loadcode = "/loadcode"
const run = "/run"

let memory = {}
let ready = false
let p = new Function("units", "enemy_units", "dropzone", "enemy_dropzone", "flags",
    "enemy_flags", "obstacles", "projectiles", "game", "memory", "{}");

const http = require('http');
const {
    parse
} = require('querystring');

const server = http.createServer((req, res) => {
    console.log(req.url)
    if (req.method === 'POST' && req.url === loadcode) {
        let body = [];
        req.on('data', (chunk) => {
            body.push(chunk);
        }).on('end', () => {
            body = Buffer.concat(body).toString();
            let params = JSON.parse(body)

            let errorText = ""
            try {
                p = new Function("units", "enemy_units", "dropzone", "enemy_dropzone", "flags",
                    "enemy_flags", "obstacles", "projectiles", "game", "memory", params.code);
            } catch (err) {
                errorText = `ERROR: ${err.name}; ${err.message};`
            }
            let resp = JSON.stringify({
                "error": errorText,
            })
            res.setHeader('Content-Type', 'application/json');
            res.statusCode = 200;
            res.end(resp)
        });
    } else if (req.method === 'POST' && req.url === run) {
        let body = [];
        req.on('data', (chunk) => {
            body.push(chunk);
            console.log(chunk);
        }).on('end', () => {
            body = Buffer.concat(body).toString();
            let params = JSON.parse(body)
            let units = params.units
            let enemy_units = params.enemy_units
            let dropzone = params.dropzone
            let enemy_dropzone = params.enemy_dropzone
            let flags = params.flags
            let enemy_flags = params.enemy_flags
            let obstacles = params.obstacles
            let projectiles = params.projectiles
            let game = params.game


            let cnsl = {}
            cnsl.logs = []
            cnsl.log = function (msg) {
                this.logs.push(msg)
            }


            let errorText = ""
            try {
                p(
                    units,
                    enemy_units,
                    dropzone,
                    enemy_dropzone,
                    flags,
                    enemy_flags,
                    obstacles,
                    projectiles,
                    game,
                    cnsl, memory)
            } catch (err) {
                errorText = `[TICK: ${10000 - params.game.ticks_left}] ERROR: ${err.name}; ${err.message};`
                cnsl.logs.push(`[TICK: ${10000 - params.game.ticks_left}] ERROR: ${err.name}; ${err.message};`)
            }

            let resp = JSON.stringify({
                "units": units,
                "enemy_units": enemy_units,
                "dropzone": dropzone,
                "enemy_dropzone": enemy_dropzone,
                "flags": flags,
                "enemy_flags": enemy_flags,
                "obstacles": obstacles,
                "projectiles": projectiles,
                "game": game,
                "console": cnsl.logs,
                "error": errorText
            })
            res.setHeader('Content-Type', 'application/json');
            res.end(resp)
        });
    }
});
server.listen(port);

function collectRequestData(request, callback) {
    const FORM_URLENCODED = 'application/x-www-form-urlencoded';
    if (request.headers['content-type'] === FORM_URLENCODED) {
        let body = '';
        request.on('data', chunk => {
            body += chunk.toString();
        });
        request.on('end', () => {
            callback(parse(body));
        });
    } else {
        callback(null);
    }
}