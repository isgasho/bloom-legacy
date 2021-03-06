<template>
  <v-container fill-height fluid class="pa-0">
    <v-col cols="4" lg="3" class="pa-0 blm-left-col">
      <v-toolbar elevation="0">
        <blm-group-selector />
        <v-tooltip bottom>
          <template v-slot:activator="{ on }">
            <v-btn icon v-on="on" to="/notes" exact v-if="archive">
              <v-icon>mdi-note</v-icon>
            </v-btn>
            <v-btn icon v-on="on" to="/notes/archive" exact v-else>
              <v-icon>mdi-archive</v-icon>
            </v-btn>
          </template>
          <span v-if="archive">Go to Notes</span>
          <span v-else>Go to Archive</span>
        </v-tooltip>
        <v-spacer />
        <v-tooltip bottom v-if="!archive">
          <template v-slot:activator="{ on }">
            <v-btn icon v-on="on" @click="newNote">
              <v-icon>mdi-pencil-plus</v-icon>
            </v-btn>
          </template>
          <span>New Note</span>
        </v-tooltip>
      </v-toolbar>

      <div style="height: calc(100vh - 65px)" class="overflow-y-auto">
        <v-alert icon="mdi-alert-circle" type="error" dismissible :value="error !== ''">
          {{ error }}
        </v-alert>
        <v-list-item-group
          v-model="selectedNoteIndex"
          @change="setSelectedNoteIndex"
          color="indigo">
          <v-list three-line class="pa-0">
            <template v-for="(note, index) in notes" class="blm-pointer">
              <v-list-item :key="`note-${index}`">

                <v-list-item-content class="text-left">
                  <v-list-item-title>{{ note.data.title }}</v-list-item-title>
                  <v-list-item-subtitle>{{ note.data.body }}</v-list-item-subtitle>
                </v-list-item-content>

              </v-list-item>
              <v-divider v-if="index !== notes.length - 1" :key="index"/>
            </template>
          </v-list>
        </v-list-item-group>
      </div>
    </v-col>

    <v-col cols="8" lg="9" class="pa-0">
      <blm-notes-note v-if="selectedNote"
        :note="selectedNote"
        @archived="noteArchived"
        @unarchived="noteUnarchived"
        @deleted="noteDeleted"
      />
    </v-col>
  </v-container>
</template>


<script lang="ts">
import {
  Component, Prop, Vue, Watch,
} from 'vue-property-decorator';
import BlmNote from './note.vue';
import core, { BlmObject, Notes } from '@/core';
import {
  Note,
  Method,
  NOTE_TYPE,
} from '@/core/notes';
import { log } from '@/libs/rz';
import BlmGroupSelector from '@/ui/kernel/components/group_selector.vue';
import { NotesFindParams, NotesCreateParams } from '@/core/messages';
import { Group } from '@/api/models';


@Component({
  components: {
    'blm-notes-note': BlmNote,
    BlmGroupSelector,
  },
})
export default class BlmNotes extends Vue {
  // props
  @Prop({ type: Boolean, default: false }) archive!: boolean;

  // data
  error = '';
  isLoading = false;
  notes: BlmObject<Note>[] = [];
  selectedNote: BlmObject<Note> | null = null;
  selectedNoteIndex: number | undefined = 0;
  saveInterval: any | null = null;
  lastSavedTitle = '';
  lastSavedBody = '';

  // computed
  // lifecycle
  async created() {
    await this.load(this.$store.state.selectedGroup);
    this.saveInterval = setInterval(this.save, 2000);
  }

  async beforeDestroy() {
    await this.save();
  }

  destroyed() {
    if (this.saveInterval) {
      clearInterval(this.saveInterval);
    }
  }

  // watch
  @Watch('$store.state.selectedGroup')
  async onSelectedGroupChanged(group: Group | null) {
    await this.save();
    this.load(group);
  }

  // methods
  async load(group: Group | null) {
    if (this.archive) {
      await this.findArchived(group);
    } else {
      await this.findNotes(group);
    }
    this.setSelectedNoteIndex(0);
  }

  async save() {
    if (this.selectedNote) {
      if (this.selectedNote.id === '') {
        await this.createNote();
      } else {
        await this.updateNote();
      }
    }
  }

  async findNotes(group: Group | null) {
    this.error = '';
    this.isLoading = true;
    const params: NotesFindParams = {
      groupID: group?.id || null,
    };

    try {
      const res = await core.call(Method.FindNotes, params);
      this.notes = (res as Notes).notes;
    } catch (err) {
      log.error(err);
    } finally {
      this.isLoading = false;
    }
  }

  async findArchived(group: Group | null) {
    this.error = '';
    this.isLoading = true;
    const params: NotesFindParams = {
      groupID: group?.id || null,
    };

    try {
      const res = await core.call(Method.FindNotes, params);
      this.notes = (res as Notes).notes;
    } catch (err) {
      this.error = err.message;
    } finally {
      this.isLoading = false;
    }
  }

  // noteCreated(note: Note) {
  //   this.notes = [note, ...this.notes];
  // }

  noteUpdated(updatedNote: BlmObject<Note>) {
    // const pos = this.notes.map((note: Note) => note.id).indexOf(updatedNote.id);
    // this.notes.splice(pos, 1);
    // this.notes = [updatedNote, ...this.notes];
    this.notes[this.selectedNoteIndex!] = updatedNote;
    this.selectedNote = updatedNote;


    // this.notes = this.notes.map((note: any) => {
    //   if (note.id === updatedNote.id) {
    //     return updatedNote;
    //   }
    //   return note;
    // });
  }

  noteArchived(archivedNote: BlmObject<Note>) {
    this.notes = this.notes.filter((note: BlmObject<Note>) => note.id !== archivedNote.id);
    this.selectedNote = null;
    this.setSelectedNoteIndex(0);
  }

  noteUnarchived(unarchivedNote: BlmObject<Note>) {
    this.notes = this.notes.filter((note: BlmObject<Note>) => note.id !== unarchivedNote.id);
    this.selectedNote = null;
    this.setSelectedNoteIndex(0);
  }

  noteDeleted(deletedNote: BlmObject<Note>) {
    this.notes = this.notes.filter((note: BlmObject<Note>) => note.id !== deletedNote.id);
    this.selectedNote = null;
    this.setSelectedNoteIndex(0);
  }

  async setSelectedNoteIndex(selected: number | undefined) {
    // save before changing / closing note
    await this.save();

    if (selected === undefined || selected >= this.notes.length) {
      this.selectedNoteIndex = undefined;
      this.selectedNote = null;
    } else {
      const note = this.notes[selected];
      this.notes.splice(selected, 1);
      this.selectedNote = note;
      this.notes = [note, ...this.notes];
      this.selectedNoteIndex = 0;
      this.lastSavedTitle = this.selectedNote.data.title;
      this.lastSavedBody = this.selectedNote.data.body;
    }
  }

  async newNote() {
    await this.setSelectedNoteIndex(undefined);

    const newNote: BlmObject<Note> = {
      id: '',
      createdAt: new Date(),
      updatedAt: new Date(),
      data: {
        title: '',
        body: '',
        color: '#ffffff',
        archivedAt: null,
        isFavorite: false,
      },
      groupID: null,
      type: NOTE_TYPE,
    };
    this.notes = [newNote, ...this.notes];
    await this.setSelectedNoteIndex(0);
  }

  async createNote() {
    this.error = '';
    this.isLoading = true;
    const params: NotesCreateParams = {
      title: this.selectedNote!.data.title,
      body: this.selectedNote!.data.body,
      color: '#ffffff',
      groupID: this.$store.state.selectedGroup?.id || null,
    };
    try {
      const res = await core.call(Method.CreateNote, params);
      this.notes[0] = res;
      this.selectedNote = res;
      // this.selectedNote = res;
      // this.$emit('created', (res as Note));
    } catch (err) {
      this.error = err.message;
    } finally {
      this.isLoading = false;
    }
  }

  async updateNote() {
    if (this.lastSavedTitle === this.selectedNote!.data.title
      && this.lastSavedBody === this.selectedNote!.data.body) {
      return;
    }

    this.error = '';
    this.isLoading = true;
    const note = { ...this.selectedNote } as BlmObject<Note>;
    try {
      const res = await core.call(Method.UpdateNote, note);
      this.notes[0] = res;
      this.selectedNote = res;

      this.lastSavedBody = this.selectedNote!.data.body;
      this.lastSavedTitle = this.selectedNote!.data.title;
    } catch (err) {
      this.error = err.message;
    } finally {
      this.isLoading = false;
    }
  }
}
</script>


<style lang="scss" scoped>
.blm-left-col {
  border-right: 1px solid #dedede;
}

.v-overflow-btn .v-input__slot::before {
    border-color: grey !important;
}

.v-toolbar {
  border-bottom: 1px solid rgba($color: #000000, $alpha: 0.1) !important;
  left: 0px !important;
}
</style>
