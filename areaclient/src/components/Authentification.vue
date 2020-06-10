<template>
    <div class="Auth">
        <h1>Client AreaGO</h1>
        <p>Login obbligatoire lol</p>
        <input v-model="user" placeholder="Username">
        <input v-model="passwd" placeholder="Passwd">
        <button href="#" v-on:click="Login()">Send</button>
        <p>Pas de compte ? Register ici</p>
        <input v-model="newUser" placeholder="Username">
        <input v-model="newPasswd" placeholder="Passwd">
        <button href="#" v-on:click="Register()">Send</button>
        <div>
          <button href="#" v-on:click="Ping()">Ping</button>
        </div>
        <div>
          <button href="#" v-on:click="ping()">ping</button>
        </div>
    </div>
</template>

<script>
import Vue from 'vue'
import VueResource from 'vue-resource'
Vue.use(VueResource)
Vue.http.options.emulateJSON = true
const http=Vue.http

export default {
  name: 'Auth',


    data() {
        return {
            user: '',
            passwd: '',
            newUser: '',
            newPasswd: '',
            token: ''
            }
    },

  components: {
  },

  methods: {
    Login: function () {
        //faire la requête send au serveur et 
        // voir comment lui envoyer des params et les lire dessus
        console.log("event Login avec comme param " + this.user + " " + this.passwd)
        this.$http.post('http://localhost:6060/login', {user: this.user, passwd: this.passwd}).then(function(data){
          //console.log(data)
          this.token = data.body;
          console.log("token = " + this.token)
      })
    },

    Register: function() {
        console.log("event Register avec comme param " + this.newUser + " " + this.newPasswd)
        this.$http.post('http://localhost:6060/register', {newUser: this.newUser, newPasswd: this.newPasswd}).then(function(data){
          console.log(data)
      })
    },

    /*
    Ping: function() {
      this.$http.get('http://localhost:6060/ping', {headers: {
        Authorization: this.token}
    }).then(function(data){
        console.log(this.token)
        console.log(data)
      })
    }
    */
   
      Ping: function() {
      this.$http.get('http://localhost:6060/Ping').then(function(data){
        console.log(data)
      })
    },

      ping: function() {
      this.$http.get('http://localhost:6060/ping', {headers: {'Authorization': this.token}}).then(function(data){
        console.log(data)
      })
    }
  }
}
</script>

<style>

</style>
