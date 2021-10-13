<template>
  <v-form
    ref="form"
    v-model="form.valid"
    :disabled="form.loading"
    method="POST"
    @submit.prevent="doSubmit"
  >
    <v-card>
      <v-card-title>
        Create Page
      </v-card-title>

      <v-card-text>
        <v-row>
          <v-col cols="12">
            <v-text-field
              v-model="form.inputs.name"
              :rules="form.rules.name"
              label="Name"
              required
            />
          </v-col>

          <v-col
            cols="12"
            md="6"
          >
            <v-file-input
              v-model="form.inputs.background"
              accept="image/*"
              clearable
              label="Background"
            />
          </v-col>

          <v-col
            cols="12"
            md="6"
          >
            <v-file-input
              v-model="form.inputs.backgroundMobile"
              accept="image/*"
              clearable
              label="Background Mobile"
            />
          </v-col>

          <v-col cols="12">
            <div>
              <client-only>
                <markdown-editor
                  ref="editor"
                  :initial-value="form.inputs.markdown"
                  :options="editorOptions"
                  initial-edit-type="wysiwyg"
                  @change="onEditorChange"
                />
              </client-only>
            </div>
          </v-col>
        </v-row>

        <v-row>
          <v-col cols="12">
            <iframe
              :src="$router.resolve({
                path: '/api/page/rendered',
                query: {
                  p: form.inputs.markdown,
                },
              }).href.replace(/^\/admin\//, '/')"
              style="width: 100%; height: 500px;"
            />
          </v-col>
        </v-row>
      </v-card-text>

      <v-card-actions>
        <v-spacer />

        <v-btn
          :disabled="!form.valid"
          :loading="form.loading"
          color="success"
          large
          text
          type="submit"
        >
          Save
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-form>
</template>

<router>
name: PagePagesCreate
</router>

<script>
  import {
    mapActions,
  } from "vuex";

  import "@toast-ui/editor/dist/toastui-editor.css";
  import "@toast-ui/editor/dist/theme/toastui-editor-dark.css";

  export default {
    components: {
      MarkdownEditor: () =>
        process.client
          ? import("@toast-ui/vue-editor").then(({ Editor }) => Editor)
          : () => ({
            render: (h) => h("div"),
          })
      ,
    },

    data: () => ({
      editorOptions: {
        language: "hr-HR",
        hideModeSwitch: true,
        usageStatistics: false,
        toolbarItems: [
          [ "heading" ],
          [ "bold", "italic", "strike" ],
          [ "hr", "quote" ],
        ],
        theme: "dark",
      },
      form: {
        valid: false,
        inputs: {
          name: "",
          markdown: "",
          background: null,
          backgroundMobile: null,
        },
        rules: {
          username: [
            (v) => !!v || "Name required",
          ],
          markdown: [
            (v) => !!v || "Markdown required",
          ],
        },
        loading: false,
      },
    }),

    methods: {
      ...mapActions("page", [
        "create",
      ]),

      onEditorChange() {
        this.form.inputs.markdown = this.$refs.editor.invoke("getMarkdown");
      },

      async doSubmit() {
        const { inputs } = this.form;

        this.form.loading = true;
        try {
          const res = await this.create(inputs);

          if ("success" === res.status) {
            return this.$router.push({
              name: "PagePagesList",
            });
          } else {
            alert(res.message);
          }
        } catch (e) {
          console.error(e);
          alert("Something went wrong");
        } finally {
          this.form.loading = false;
        }
      },
    },
  };
</script>
