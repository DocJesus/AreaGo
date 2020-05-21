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
            newPasswd: ''
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
          console.log(data)
      })
    },

    Register: function() {
        console.log("event Register avec comme param " + this.newUser + " " + this.newPasswd)
        this.$http.post('http://localhost:6060/register', {newUser: this.newUser, newPasswd: this.newPasswd}).then(function(data){
          console.log(data)
      })
  }
  
  }
}
</script>

<style>

</style>
