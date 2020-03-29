const stompit = require('stompit');
const connectionManager = new stompit.ConnectFailover();

connectionManager.addServer({
    'host': 'localhost',
    'port': 61613,
    'connectHeaders':{
        'host': '/',
        'login': 'username',
        'passcode': 'password',
        'heart-beat': '5000,5000'
    }
});

const channel = new stompit.Channel(connectionManager);

let temperatura = Math.floor(Math.random() * 50);
let iluminacion = Math.floor(Math.random() * (1000 - 200)) + 200;

channel.subscribe({ destination: '/topic/ActuadorTemperatura2' }, (err, msg) => {
    msg.readString('UTF-8', (err, body) => {
        console.log("Subiendo la temperatura a " +body);
        temperatura = parseInt(body);
    });
});

channel.subscribe({destination: '/topic/ActuadorIluminacion2'}, (err,msg) => {
    msg.readString('UTF-8', (err, body) => {
        console.log("Subiendo la iluminación a " + body);
        iluminacion = parseInt(body);
    });
});


setInterval(() =>{
    const sendHeaders = {
        'destination': '/topic/LecturasTemperaturas2',
        'content-type': 'text/plain'
    };

    channel.send(sendHeaders, '' + temperatura);
    temperatura = Math.floor(Math.random() * 50);
},5000);

setInterval(() => {
    const sendHeaders = {
        'destination': '/topic/LecturasIluminacion2',
        'content-type': 'text/plain'
    };

    channel.send(sendHeaders, '' + iluminacion);
    iluminacion = Math.floor(Math.random() * 50);
}, 5000);
