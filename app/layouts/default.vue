<template>
  <div class="bg-light">
    <div id="db-wrapper" class="toggled">
      <!-- <SideNavbar></SideNavbar> -->
      <div id="page-content">
        <Header></Header>
        <div class="container-fluid pt-6">
          <slot />
        </div>
      </div>
    </div>
  </div>
</template>
<script setup>
import Header from "~~/components/Header.vue";
import SideNavbar from "~~/components/SideNavbar.vue";
</script>
