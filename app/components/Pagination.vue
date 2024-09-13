<template>
    <div class="row mt-3">
        <div class="col-sm-12 col-md-6">
            <nav aria-label="Page navigation example">
                <ul class="pagination pagination-primary justify-content-end">
                    <li class="page-item">
                        <NuxtLink class="page-link" :class="{ disabled: props.page <= 1 }"
                            :to="{ path: route.path, query: { ...route.query, page: props.page - 1 } }">
                            <span aria-hidden="true">
                                <font-awesome-icon icon="caret-left" />
                            </span>
                        </NuxtLink>
                    </li>

                    <li class="page-item active">
                        <div class="page-link border-circle">{{ props.page }}</div>
                    </li>

                    <li class="page-item">
                        <NuxtLink class="page-link" :class="{ disabled: props.numOfRecords < 10 }"
                            :to="{ path: route.path, query: { ...route.query, page: props.page + 1 } }">
                            <span aria-hidden="true">
                                <font-awesome-icon icon="caret-right" />
                            </span>
                        </NuxtLink>
                    </li>
                </ul>
            </nav>
        </div>
    </div>
</template>

<script setup>
const route = useRoute();
const props = defineProps([ 'page', 'numOfRecords']); 
</script>