<template>
  <div class="card">
    <h2>Clusters</h2>
    <ul v-if="clusters.length > 0">
      <li v-for="name in clusters" :key="name">
        <span>{{ name }}</span>
        <button @click="confirmDelete(name)" :disabled="loading === name">
          {{ loading === name ? 'Deleting...' : 'Delete' }}
        </button>
      </li>
    </ul>
    <p v-else style="color:#888">No clusters found.</p>
    <button class="refresh" @click="fetchClusters">Refresh</button>
    <p v-if="message" :class="messageType">{{ message }}</p>
  </div>
</template>

<script setup>
import { ref, defineExpose } from 'vue'
const clusters = ref([])
const loading = ref('')
const message = ref('')
const messageType = ref('')

defineExpose({ fetchClusters })

async function fetchClusters() {
  loading.value = ''
  message.value = ''
  try {
    const res = await fetch('/api/list')
    if (!res.ok) throw new Error('Error listing clusters')
    clusters.value = await res.json()
  } catch (err) {
    message.value = err.message
    messageType.value = 'error'
  }
}
fetchClusters()

const emit = defineEmits(['deleted'])
async function confirmDelete(name) {
  if (!confirm(`Delete cluster "${name}"? This cannot be undone.`)) return
  loading.value = name
  message.value = ''
  try {
    const res = await fetch(`/api/delete?name=${encodeURIComponent(name)}`, { method: 'DELETE' })
    if (!res.ok) throw new Error(await res.text())
    clusters.value = clusters.value.filter(n => n !== name)
    message.value = 'Cluster deleted.'
    messageType.value = 'success'
    emit('deleted')
  } catch (err) {
    message.value = err.message
    messageType.value = 'error'
  } finally {
    loading.value = ''
    setTimeout(() => (message.value = ''), 2000)
  }
}
</script>

<style scoped>
.card {
  margin-bottom: 2rem;
  background: #fff;
  border-radius: 12px;
  box-shadow: 0 1px 6px #0001;
  padding: 1.5rem;
}
ul {
  list-style: none;
  padding: 0;
  margin: 0 0 1rem 0;
}
li {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0.4rem 0;
  border-bottom: 1px solid #f1f1f1;
}
li:last-child {
  border-bottom: none;
}
button {
  background: #e74c3c;
  color: #fff;
  border: none;
  border-radius: 5px;
  padding: 0.2rem 0.8rem;
  cursor: pointer;
  transition: background 0.2s;
}
button:disabled {
  background: #f1adad;
}
.refresh {
  margin-top: 1rem;
  background: #4f8cfb;
}
.success {
  color: #258b40;
  margin-top: 0.5rem;
}
.error {
  color: #e25241;
  margin-top: 0.5rem;
}
</style>
