<template>
  <div class="hello">
    <!-- <h1>{{ msg }}</h1> -->
    
    <div>
      <h1 class="title">Buscador De Dominios</h1>
    </div>
    <hr>
    
    <div>
      <b-container fluid>
        <b-row class="text-center" align-h="center">       
          <form>
            <div class="field">
              <label class="label">Ingrese la URL del dominio  ej: facebook.com (solo dominio.com)</label>
              <b-form-input name="domain" v-model="domain"  class="input" type="text"></b-form-input>
                <div class="mt-2">
                  <b-button-group vertical>
                    <b-button v-b-toggle.collapse-1 variant="success" v-on:click="endpoint1()">Buscar</b-button>
                      <b-collapse id="collapse-1" class="mt-2">
                        <b-overlay :show="show" rounded="sm">
                          <b-card title="Infomación del dominio" :aria-hidden="show ? 'true' : null">
                             
                                <b-col cols="20" class="text-left">
                                  <div><label class="label">Title : {{ title }}</label></div>
                                  <div><label class="label">Logo: {{ logo }}</label></div> 
                                  
                                  <b-img v-bind:src="logo" rounded="circle" alt="Circle image"  contain  height="100px"  width="150px"></b-img>
                                  <div><label class="label">Is down: {{ isdown }}</label></div>          
                                  <div><label class="label">SSL grade: {{ sslgrade }}</label></div>   
                                  <div><label class="label">Previous ssl grade: {{ prevgrade }}</label></div>          
                                  <div><label class="label">Have servers changed: {{ serverschange }}</label></div>                  
                                </b-col>
                                <b-button v-b-toggle.collapse-1-inner size="sm">Información de los servidores</b-button>
                                  <b-collapse id="collapse-1-inner" class="mt-2">
                                    <b-card>
                                      <table class="table">
                                        <thead>
                                          <tr>
                                            <th scope="col">#</th>
                                            <th scope="col">address</th>
                                            <th scope="col">SSL grade</th>
                                            <th scope="col">Country</th>
                                            <th scope="col">Owner</th>
                                          </tr>
                                        </thead>
                                          <tbody v-for="(server, index) in servers" :key="server.id"> 
                                            <tr>
                                              <th scope="row">{{index+1}}</th>
                                              <td>{{server.addres}}</td>
                                              <td>{{server.ssl_grade}}</td>
                                              <td>{{server.country}}</td>
                                              <td>{{server.owner}}</td>
                                            </tr>                                  
                                          </tbody>
                                      </table>
                                    </b-card>
                                  </b-collapse>
                          </b-card>
                        </b-overlay>
                      </b-collapse>
                     
                    <b-button v-b-toggle.collapse-2 variant="info" v-on:click="endpoint2()">Historial</b-button>
                      <b-collapse id="collapse-2" class="mt-2">
                        <b-overlay :show="mostrar" rounded="sm">
                          <b-card title="Dominios buscados" :aria-hidden="mostrar ? 'true' : null">
                        
                               <p class="card-text">Dominios</p>
                                <b-col cols="20" class="text-left">
                                <div> <b-list-group-item>Buscados : {{ items }}</b-list-group-item> </div>
                                </b-col>
                              
                          </b-card>
                       </b-overlay>
                     </b-collapse> 
                  </b-button-group>
                </div>
              <label class="label">En Historial puedes ver los dominios buscados previamente</label>
            </div>
          </form>
        </b-row>
      </b-container>
    </div>

    &nbsp;
    

    &nbsp;
    <hr>

  </div>
</template>

<script>


import axios from 'axios';
import Vue from 'vue'
import * as VeeValidate from 'vee-validate'

/* eslint-disable */
Vue.use(VeeValidate)

export default {
  name: 'HelloWorld',  


  data: function() {
    return {
      servers: [], serverschange : "", sslgrade:"", prevgrade:"", logo:"",title:"",isdown:"",
      domain: "", items:"" , show: true, mostrar:true
    }
  },

  methods: {
    endpoint1: function() {
      var data = {"domain": this.domain}

      /*eslint-disable*/
      console.log(data) 
      /*eslint-enable*/

      axios({ method: "POST", url: "http://localhost:8000/buscar", data: data, headers: {"content-type": "text/plain" } }).then(result => { 
          // this.response = result.data;
          this.title =result.data['title']
          this.logo = result.data['logo']
          this.isdown = result.data['is_down']
          this.prevgrade = result.data['previous_ssl_grade']
          this.sslgrade = result.data['ssl_grade']
          this.serverschange = result.data['servers_changed']
          this.servers = result.data['servers']
          this.show =false
          /*eslint-disable*/
          console.log(result.data) 
          /*eslint-enable*/

        }).catch( error => {
            /*eslint-disable*/
            console.error(error);
            /*eslint-enable*/
      });
    }, 
    endpoint2: function(){
      axios({ method: "GET", url: "http://localhost:8000/historial", headers : {"content-type": "text/plain"}}).then(result => {
          this.items = result.data['items']
          this.mostrar =false
          console.log(result) 
        }).catch( error => {
            /*eslint-disable*/
            console.error(error);
            /*eslint-enable*/
      });
    }
  },
  const: {
  transformAssetUrls: {
    video: ['src', 'poster'],
    source: 'src',
    img: 'src',
    image: 'xlink:href',
    'b-avatar': 'src',
    'b-img': 'src',
    'b-img-lazy': ['src', 'blank-src'],
    'b-card': 'img-src',
    'b-card-img': 'src',
    'b-card-img-lazy': ['src', 'blank-src'],
    'b-carousel-slide': 'img-src',
    'b-embed': 'src'
  }
}
  
}

</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
h3 {
  margin: 40px 0 0;
}
ul {
  list-style-type: none;
  padding: 0;
}
li {
  display: inline-block;
  margin: 0 10px;
}
a {
  color: #42b983;
}
</style>