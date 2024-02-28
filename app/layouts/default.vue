<template>
  <div>
    <header class="container-fluid ps-0 pe-0">
      <nav class="navbar navbar-expand-lg navbar-dark bg-dark">
        <div class="container-fluid p-2">
          <!-- Title on the left -->
          <NuxtLink class="navbar-brand" href="/">Quiz App</NuxtLink>

          <!-- Profile/Login on the right -->
          <div class="navbar-collapse justify-content-end">
            <ul class="navbar-nav">
              <li class="nav-item">
                <NuxtLink class="nav-link" href="/account/login">Profile/Login</NuxtLink>
              </li>
            </ul>
          </div>
        </div>
      </nav>
    </header>
    <main class="container">
      <slot />
    </main>
    <footer
      class="container-fluid position-fixed bottom-0 bg-dark text-capitalize fs-5 p-2"
    >
      <div class="text-center p-3 text-light">
        © {{ new Date().getFullYear() }} Copyright:
        <NuxtLink class="text-light" href="/">Quizz_app</NuxtLink>
      </div>
    </footer>
    <footer class="bg-body-tertiary text-center text-lg-start"></footer>
  </div>
</template>
