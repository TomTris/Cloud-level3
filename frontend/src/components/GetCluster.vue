<template>
  <div class="card">
    <h2>Get Cluster Info</h2>
    <form @submit.prevent="getCluster">
      <input v-model="name" placeholder="Cluster Name" required />
      <button :disabled="loading">{{ loading ? "Loading..." : "Get" }}</button>
    </form>
    <pre v-if="info" class="json">{{ info }}</pre>
    <p v-if="message" :class="messageType">{{ message }}</p>
  </div>
</template>

<script setup>
import { ref } from 'vue'
const name = ref('')
const info = ref('')
const loading = ref(false)
const message = ref('')
const messageType = ref('')

async function getCluster() {
  loading.value = true
  info.value = ''
  message.value = ''
  try {
    const res = await fetch(`/api/get?name=${encodeURIComponent(name.value)}`)
    if (!res.ok) throw new Error(await res.text())
    const data = await res.json()
    info.value = JSON.stringify(data, null, 2)
    message.value = ''
  } catch (err) {
    message.value = err.message
    messageType.value = 'error'
  } finally {
    loading.value = false
    setTimeout(() => (message.value = ''), 2500)
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
form {
  display: flex;
  gap: 0.5rem;
  margin-bottom: 0.7rem;
}
input {
  flex: 1;
  padding: 0.5rem;
  border-radius: 6px;
  border: 1px solid #ccc;
  font-size: 1rem;
}
button {
  background: #4f8cfb;
  color: #fff;
  border: none;
  border-radius: 6px;
  padding: 0.5rem 1rem;
  font-size: 1rem;
  cursor: pointer;
  transition: background 0.2s;
}
button:disabled {
  background: #88bffc;
  cursor: not-allowed;
}
.json {
  background: #f7f8fa;
  border-radius: 8px;
  padding: 1rem;
  font-size: 0.93rem;
  color: #111;
  max-width: 100%;
  overflow-x: auto;
}
.error {
  color: #e25241;
  margin-top: 0.3rem;
}
</style>
