<template>
    <section class="flex flex-col h-full">
        <modal
            title="New Project"
            description="Please fill in a name for your new project"
            actionButtonText="Create Project"
            v-if="showNewProjectModal"
            @close="showNewProjectModal = false"
            @action="createNewProject"
            :enableOK="nameValidationErrors(newProjectName).length === 0"
        >
            <input type="text" placeholder="chapi.com external api's" class="w-full" v-model="newProjectName" @input="isProjectCreated">
            <template v-if="nameValidationErrors(newProjectName).length > 0">
                <p class="text-sm error-text" v-for="(error, index) in nameValidationErrors(newProjectName)" :key="index">{{ error }}</p>
            </template>
        </modal>

        <modal
            title="Update Project Name"
            description="Enter new name of project"
            actionButtonText="Update Project"
            v-if="showEditModal"
            @close="showEditModal = false"
            @action="updateProject"
            :enableOK="nameValidationErrors(editProjectName).length === 0"
        >
            <input type="text" :placeholder="projectToEdit.name" class="w-full" v-model="editProjectName" @input="isProjectCreated">
            <template v-if="nameValidationErrors(editProjectName).length > 0">
                <p class="text-sm error-text" v-for="(error, index) in nameValidationErrors(editProjectName)" :key="index">{{ error }}</p>
            </template>
        </modal>

        <div class="px-6 py-8 border-b border-gray-200">
            <input 
                type="text" 
                placeholder="Search Projects" 
                v-model="projectSearchText"
                class="w-full"
            >
        </div>

        <div class="px-6 py-4">
            <button class="w-full" @click="showNewProjectModal = true">Create New Project</button>
        </div>

        <ul class="pt-4 px-4 overflow-y-auto h-auto">
            <li 
                v-for="project in filteredProjects" 
                class="px-4 py-2 flex items-center cursor-pointer"
                :class="{'active': selectedProject.id === project.id}"
                :key="project.id"
                @click="getProjectRoutes(project)"
            >
                <span>{{ project.name }}</span>
                <span class="inline-block text-xl ml-auto text-gray-200 hover:text-red-500 transition duration-300" @click="deleteProject(project.id)">
                    <i class="ri-delete-bin-line"></i>
                </span>
                <span class="inline-block text-xl ml-6 text-gray-200 hover:text-blue-500 transition duration-300" @click="initProjectNameEdit(project)">
                    <i class="ri-edit-line"></i>
                </span>
            </li>
        </ul>
    </section>
</template>

<script lang='ts'>
import {Vue, Component, Watch} from 'vue-property-decorator';

import Modal from '@/components/Modal/Modal.vue';

import { projects, fetchUserProjects, createProject, deleteProject, updateProject, isProjectCreated } from '@/store/modules/project';
import { authenticatedUser } from '@/store/modules/user';
import { getRoutes } from '@/store/modules/route';

import User from '@/types/User';
import Project from '@/types/Project';
import Route from '@/types/Route';
import { Route as VueRoute } from 'vue-router';

@Component({
    components: {
        Modal
    }
})
export default class ProjectNav extends Vue {
    private newProjectName = '';

    private selectedProject: Project = new Project;

    private showNewProjectModal = false;

    private projectSearchText = '';

    // set by the isProjectCreated method further below
    private projectAlreadyExists = false;

    private protectedProjectNames: string[] = ['www', 'chapi', 'localhost'];
    private nameValidationErrors(name: string): string[] {
        let errors = [];

        if (name.length < 3) errors.push('Project name must have at least 3 characters');
        if (this.projectAlreadyExists) errors.push('Project already exists');
        if (this.protectedProjectNames.includes(name.toLowerCase())) errors.push(`${name} is not allowed as a valid project name`);
        if (/^[A-Za-z0-9-]*$/.test(name) === false) errors.push('A project name can only have letters, numbers and hyphen\n');

        return errors
    }

    get user(): User | null {
        return authenticatedUser(this.$store)
    }

    get projects(): Project[] {
        return projects(this.$store)
    }

    get filteredProjects(): Project[] {
        if (this.projectSearchText.length === 0) return this.projects;

        return this.projects.filter(project => {
            return project.name
                .toLowerCase()
                .includes(this.projectSearchText.toLowerCase())
        })
    }

    get projectRoutes(): Route[] {
        return getRoutes(this.$store)
    }

    private viewFirstProject() {
        if (this.projects.length > 0) {
            this.getProjectRoutes(this.projects[0])
        }
    }

    private createNewProject() {
        if (this.nameValidationErrors(this.newProjectName).length > 0) return;

        if (this.user) {
            createProject(this.$store, {
                name: this.newProjectName,
                userId: this.user.id
            })
            .then(() => {
                this.showNewProjectModal = false;
                this.newProjectName = '';
            })
            .catch(() => this.$toast.error('Error creating project'));
        }
    }


    private editProjectName = '';
    private showEditModal = false;
    private projectToEdit = new Project();
    private initProjectNameEdit(project: Project) {
        this.editProjectName = '';
        this.showEditModal = true;
        Object.assign(this.projectToEdit, project);
    }

    private deleteProject(projectId: number) {
        deleteProject(this.$store, projectId)
        .catch(() => this.$toast.error('Error deleting project'));
    }

    private updateProject() {
        updateProject(this.$store, {
            ...this.projectToEdit,
            name: this.editProjectName
        })
        .then(() => {
            this.showEditModal = false;
        })
        .catch(() => this.$toast.error('Error updating project'));
    }

    private isProjectCreated(event: InputEvent) {
        const { value } = event.target as HTMLInputElement;

        if (value === '') return;

        isProjectCreated(this.$store, value)
            .then((doesExist) => {
                this.projectAlreadyExists = doesExist;
            })
    }

    private getProjectRoutes(project: Project) {
        this.selectedProject = project;

        this.$router.push({ 
            name: 'Routes List',
            query: {
                project: String(project.id)
            }
        })
    }

    get route(): VueRoute {
        return this.$route;
    }
    @Watch('route')
    onRouteChange(): void {
        if (this.route.path == '/dashboard') {
            this.viewFirstProject();
        }
    }

    mounted(): void {
        if (this.projects.length === 0 && this.user) {
            fetchUserProjects(this.$store)
                .then(() => {
                    this.viewFirstProject()
                })
        } else {
            this.viewFirstProject()
        }
    }
}
</script>

<style lang="postcss" scoped>
.active {
    @apply bg-gray-100 rounded-md;
}
</style>