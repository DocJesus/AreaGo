<template>
    <div class="Auth">
        <input v-model="user" placeholder="Username">
        <input v-model="passwd" placeholder="Passwd">

        <button href="#" v-on:click="Send(user, passwd)">Send</button>
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

        }
    },

  components: {
  },

  methods: {
    Send: function (user, passwd) {
        //faire la requête send au serveur et 
        // voir comment lui envoyer des params et les lire dessus
        console.log("event Send avec comme param " + user + " " + passwd)
        this.$http.post('http://localhost:8080/ping', {user: user, passwd: passwd}).then(function(data){
          console.log(data)
      })
    }
  }

}
</script>

<style>

</style>
