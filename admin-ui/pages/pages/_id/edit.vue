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
            <div
              v-if="form.inputs.backgroundId"
            >
              <v-img
                :src="page.background.srcset[page.background.srcset.length - 1].url"
                contain
                max-height="150px"
              />
              <v-btn
                color="error"
                icon
                @click="form.inputs.backgroundId = null"
              >
                <v-icon>mdi-delete</v-icon>
              </v-btn>
            </div>
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
            <div
              v-if="form.inputs.backgroundMobileId"
            >
              <v-img
                :src="page.backgroundMobile.srcset[page.backgroundMobile.srcset.length - 1].url"
                contain
                max-height="150px"
              />
              <v-btn
                color="error"
                icon
                @click="form.inputs.backgroundMobileId = null"
              >
                <v-icon>mdi-delete</v-icon>
              </v-btn>
            </div>
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
          <v-col :class="$style.iframeContainer" cols="12">
            <transition>
              <div
                v-if="preview.loading"
                :class="$style.iframeLoadingContainer"
              >
                <v-progress-circular
                  color="secondary"
                  indeterminate
                />
              </div>
            </transition>
            <iframe
              ref="previewFrame"
              style="width: 100%; height: 500px; border: 1px solid rgba(255, 255, 255, 0.7); border-radius: 4px;"
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
name: PagePagesEdit
</router>

<script>
  import {
    mapActions,
    mapGetters,
  } from "vuex";

  import "@toast-ui/editor/dist/toastui-editor.css";
  import "@toast-ui/editor/dist/theme/toastui-editor-dark.css";


  // Returns a function, that, when invoked, will only be triggered at most once
  // during a given window of time. Normally, the throttled function will run
  // as much as it can, without ever going more than once per `wait` duration;
  // but if you'd like to disable the execution on the leading edge, pass
  // `{leading: false}`. To disable execution on the trailing edge, ditto.
  function throttle(func, wait, options = {}) {
    let context, args, result;
    let timeout = null;
    let previous = 0;

    const later = () => {
      previous = false === options.leading ? 0 : Date.now();
      timeout = null;
      result = func.apply(context, args);
      if (!timeout) {
        context = null;
        args = null;
      }
    };

    return function(...argList) {
      const now = Date.now();
      if (!previous && false === options.leading) {
        previous = now;
      }
      const remaining = wait - (now - previous);
      context = this;
      args = argList;
      if (0 >= remaining || remaining > wait) {
        if (timeout) {
          clearTimeout(timeout);
          timeout = null;
        }
        previous = now;
        result = func.apply(context, args);
        if (!timeout) {
          context = null;
          args = null;
        }
      } else if (!timeout && false !== options.trailing) {
        timeout = setTimeout(later, remaining);
      }
      return result;
    };
  }

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

    data() {
      const page = this.$store.getters["page/page"];

      return ({
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
            name: String(page.name),
            markdown: String(page.contents || ""),
            background: null,
            backgroundId: page.background?.id,
            backgroundMobile: null,
            backgroundMobileId: page.backgroundMobile?.id,
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
        preview: {
          loading: false,
        },
      });
    },

    computed: {
      ...mapGetters("page", [
        "page",
      ]),
    },

    watch: {
      "form.inputs.markdown": throttle(
        function() {
          this.fetchRenderedPage();
        },
        1000,
        {
          leading: true,
          trailing: true,
        },
      ),
    },

    mounted() {
      this.fetchRenderedPage();
    },

    methods: {
      ...mapActions("page", [
        "update",
      ]),

      onEditorChange() {
        this.form.inputs.markdown = this.$refs.editor.invoke("getMarkdown");
      },

      async fetchRenderedPage() {
        const { inputs } = this.form;

        const data = new FormData();
        data.set("markdown", inputs.markdown);

        this.preview.loading = true;
        const resp = await this.$api.post("/page/rendered", data);
        this.preview.loading = false;

        const $preview = this.$refs.previewFrame;
        const doc = $preview.contentWindow.document || $preview.contentDocument;
        const $doc = doc.open();
        $doc.write(resp.data);
        $doc.close();
      },

      async doSubmit() {
        const { inputs } = this.form;

        this.form.loading = true;
        try {
          const res = await this.update({
            id: this.page.id,
            ...inputs,
          });

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

<style lang="scss" module>
  .iframeContainer {
    position: relative;

    /* stylelint-disable-next-line selector-pseudo-class-no-unknown */
    :global {
      .v-enter-active,
      .v-leave-active {
        transition: opacity 0.25s ease-out;
      }

      .v-enter-from,
      .v-leave-to {
        opacity: 0;
      }
    }

    .iframeLoadingContainer {
      position: absolute;
      top: 1.5em;
      left: 1.5em;
      padding: 0.5em;
      border-radius: 4px;
      background: rgba(0, 0, 0, 0.5);
    }
  }
</style>
