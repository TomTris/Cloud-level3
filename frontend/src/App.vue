<template>
  <div class="container">
    <h1>Postgres Cluster Manager</h1>
    <CreateCluster @created="refreshClusters" />
    <GetCluster ref="getRef" />
    <ClusterList ref="listRef" @deleted="refreshClusters" @get="openGet" />
  </div>
</template>

<script setup>
import { ref } from 'vue'
import CreateCluster from './components/CreateCluster.vue'
import GetCluster from './components/GetCluster.vue'
import ClusterList from './components/ClusterList.vue'

const listRef = ref(null)
const getRef = ref(null)

function refreshClusters() {
  listRef.value?.fetchClusters()
}
function openGet(name) {
  getRef.value?.showCluster(name)
}
</script>

<style>
.container {
  max-width: 600px;
  margin: 2rem auto;
  background: #f9f9f9;
  border-radius: 16px;
  box-shadow: 0 2px 12px #0001;
  padding: 2rem;
  font-family: system-ui, sans-serif;
}
h1 {
  text-align: center;
  margin-bottom: 1.5rem;
}
</style>
