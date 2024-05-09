<script setup></script>
<template>
  <div class="container">
    <div>
      <font-awesome-icon
        icon="fa-brands fa-github-square"
        class="icon-fontawesome"
      />
      <font-awesome-icon :icon="['fas', 'camera']" class="icon-fontawesome" />

      <font-awesome-icon icon="camera" />
    </div>
    <div class="d-flex">
      <!-- DropDown Code -->
      <div class="dropdown mt-3 me-4">
        <button
          id="dropdownMenuButton1"
          class="btn btn-secondary dropdown-toggle"
          type="button"
          data-bs-toggle="dropdown"
          aria-expanded="false"
        >
          Dropdown button
        </button>
        <ul class="dropdown-menu" aria-labelledby="dropdownMenuButton1">
          <li><a class="dropdown-item" href="#">Action</a></li>
          <li><a class="dropdown-item" href="#">Another action</a></li>
          <li><a class="dropdown-item" href="#">Something else here</a></li>
        </ul>
      </div>

      <!-- Button trigger modal -->
      <button
        type="button"
        class="btn btn-primary mt-3"
        data-bs-toggle="modal"
        data-bs-target="#exampleModal"
      >
        Launch demo modal
      </button>

      <!-- Modal -->
      <div
        id="exampleModal"
        class="modal fade"
        tabindex="-1"
        aria-labelledby="exampleModalLabel"
        aria-hidden="true"
      >
        <div class="modal-dialog">
          <div class="modal-content">
            <div class="modal-header">
              <h5 id="exampleModalLabel" class="modal-title">Modal title</h5>
              <button
                type="button"
                class="btn-close"
                data-bs-dismiss="modal"
                aria-label="Close"
              ></button>
            </div>
            <div class="modal-body">Add your Content Here</div>
            <div class="modal-footer">
              <button
                type="button"
                class="btn btn-secondary"
                data-bs-dismiss="modal"
              >
                Close
              </button>
              <button type="button" class="btn btn-primary">
                Save changes
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
<style scoped>
.icon-fontawesome {
  font-size: 32px;
  margin-right: 12px;
}
</style>
