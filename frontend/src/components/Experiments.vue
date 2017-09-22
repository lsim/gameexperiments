<template>
  <div v-on:click="foobar">
    experiments component
  </div>
</template>

<script>
  import FooBarBaz from './experiments'
  export default {
    name: 'experiments',
    components: {

    }, methods: {
      foobar() {
        FooBarBaz()
      }
    }
  }
</script>
<style lang="sass">

</style>
