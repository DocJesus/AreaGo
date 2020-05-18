<template>
  <div class="Blogs">
      <h1>all blogs</h1>
      <div v-for="blog in blogs" class="singleBLog">
          <h2>{{blog.title}}</h2>
          <article>{{blog.body}}</article>
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
  name: 'Blogs',

    data() {
        return {
            blogs: []
        }
    },

  components: {
  },

  created() {
      this.$http.post('http://localhost:8080/ping').then(function(data){
          console.log(data)
      })
            this.$http.get('http://localhost:8080/ping').then(function(data){
          console.log(data)
          this.blogs = data.body.slice(0, 10)
      })

       /*this.$http.get('https://jsonplaceholder.typicode.com/posts').then(function(data){
          console.log(data)
          this.blogs = data.body.slice(0, 10)
      })*/ }
}
</script>

<style>

</style>
