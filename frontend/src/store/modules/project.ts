import { ActionContext } from 'vuex';
import { getStoreAccessors } from "vuex-typescript";
import Project from '@/types/Project';
import StoreState from '@/store/storeState';
import { Response } from '@/types/HTTP';

type ProjectContext = ActionContext<ProjectState, StoreState>;

const API = process.env.VUE_APP_SERVER;

export type ProjectState = {
    projects: Project[];
}

const state: ProjectState = {
    projects: []
}

const project = {
    namespaced: true,
    state,
    getters: {
        getProjects(state: ProjectState): Project[] {
            return state.projects;
        },
        getProjectById(state: ProjectState): (projectId: number) => Project {
            return ((projectId: number): Project => {
                const index = state.projects.findIndex(project => project.id === projectId);

                return state.projects[index]
            })
        }
    },
    mutations: {
        setProjects(state: ProjectState, projects: Project[]): void {
            state.projects = projects;
        },
        addProject(state: ProjectState, project: Project): void {
            state.projects.push(project)
        },
        deleteProject(state: ProjectState, projectId: number): void {
            const index = state.projects.findIndex(project => project.id === projectId);

            if (index !== -1) {
                state.projects.splice(index, 1)
            }
        },
        removeProject(state: ProjectState, deleteProject: Project): void {
            const projectIndex = state.projects.findIndex(project => project.id === deleteProject.id);

            if (projectIndex !== -1) {
                state.projects.splice(projectIndex, 1)
            }
        }
    },
    actions: {
        fetchUserProjects(context: ProjectContext): Promise<Project[]> {
            return new Promise((resolve, reject) => {
                fetch(`${API}/project`, {
                    credentials: 'include',
                    mode: 'cors'
                })
                .then(response => response.json())
                .then((body: Response<Project[]>) => {
                    context.commit('setProjects', body.data);

                    resolve(body.data)
                })
                .catch(error => reject(error))
            })
        },
        createProject(context: ProjectContext, project: Project): Promise<void> {
            return new Promise((resolve, reject) => {
                fetch(`${API}/project?name=${project.name}`, { 
                        method: 'POST',
                        credentials: 'include'
                    })
                    .then((res) => res.json())
                    .then((body: Response<string>) => {
                        context.commit('addProject', {
                            id: body.data,
                            name: project.name,
                            userId: project.userId
                        })
                        resolve()
                    })
                    .catch((error) => {
                        reject(error)
                    })
            })
        },
        deleteProject(context: ProjectContext, projectId: number): Promise<void> {
            return new Promise((resolve, reject) => {
                fetch(`${API}/project/${projectId}`, {
                        method: 'DELETE'
                    })
                    .then((res) => res.json())
                    .then(() => {
                        context.commit('deleteProject', projectId)
                        resolve()
                    })
                    .catch((error) => {
                        reject(error)
                    })
            })
        },
        isProjectCreated(context: ProjectContext, projectName: string): Promise<boolean> {
            return new Promise((resolve, reject) => {
                fetch(`${API}/project/exists?name=${projectName}`, {
                        method: 'POST'
                    })
                    .then((res) => res.json())
                    .then((body: Response<boolean>) => {
                        if (body.successful) {
                            resolve(body.data)
                        } else {
                            reject(body.data)
                        }
                    })
                    .catch((error) => {
                        reject(error)
                    })
            })
        },
    }
}

const { read, dispatch } = getStoreAccessors<ProjectState, StoreState>('project');

export const projects = read(project.getters.getProjects);
export const getProjectById = read(project.getters.getProjectById);

export const fetchUserProjects = dispatch(project.actions.fetchUserProjects);
export const createProject = dispatch(project.actions.createProject);
export const deleteProject = dispatch(project.actions.deleteProject);
export const isProjectCreated = dispatch(project.actions.isProjectCreated);

export default project;