<script setup lang="ts">
import { provide, ref } from 'vue'
interface Props { defaultValue?: string; modelValue?: string }
const props = defineProps<Props>()
const emit = defineEmits<{ 'update:modelValue': [v: string] }>()
const active = ref(props.modelValue ?? props.defaultValue ?? '')
provide('tabs', { active, setActive: (v: string) => { active.value = v; emit('update:modelValue', v) } })
</script>
<template><div><slot /></div></template>
