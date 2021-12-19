<template>
    <section class="flex flex-col h-full">
        <modal
            title="New Project"
            description="Please fill in a name for your new project"
            actionButtonText="Create Project"
            v-if="showNewProjectModal"
            @close="showNewProjectModal = false"
            @action="createNewProject"
            :enableOK="isInputValid"
        >
            <input type="text" placeholder="chapi.com external api's" class="w-full" v-model="newProjectName" @input="isProjectCreated">
            <p class="text-sm error-text" v-if="isProjectNameProtected">{{ newProjectName }} can't be a project name</p>
            <p class="text-sm error-text" v-if="projectAlreadyExists">Name isn't available</p>
            <p class="text-sm error-text" v-if="newProjectName.length < 3">Name should have at least 3 letters</p>
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
                class="px-4 py-2 flex justify-between items-center cursor-pointer"
                :class="{'active': selectedProject.id === project.id}"
                :key="project.id"
                @click="getProjectRoutes(project)"
            >
                <span>{{ project.name }}</span>
                <span class="inline-block text-xl text-gray-200 hover:text-red-500 transition duration-300" @click="deleteProject(project.id)">
                    <i class="ri-delete-bin-line"></i>
                </span>
            </li>
        </ul>
    </section>
</template>

<script lang='ts'>
import {Vue, Component, Watch} from 'vue-property-decorator';

import Modal from '@/components/Modal/Modal.vue';

import { projects, fetchUserProjects, createProject, deleteProject, isProjectCreated } from '@/store/modules/project';
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

    private projectAlreadyExists = false;

    private protectedProjectNames: string[] = ['www', 'chapi', 'localhost'];
    get isProjectNameProtected(): boolean {
        return this.protectedProjectNames.includes(this.newProjectName.toLowerCase());
    }

    get isInputValid(): boolean {
        return this.newProjectName.length >= 3 && !this.projectAlreadyExists;
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
        if (!this.isInputValid) return;

        if (this.user) {
            createProject(this.$store, {
                name: this.newProjectName,
                userId: this.user.id
            })
            .then(() => {
                this.showNewProjectModal = false;
                this.newProjectName = '';
            })
        }
    }

    private deleteProject(projectId: number) {
        deleteProject(this.$store, projectId)
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