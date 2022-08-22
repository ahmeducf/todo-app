<script lang="ts">

    import Task from "./lib/Task.svelte"
    import { onMount } from "svelte"

    const URL:string = "http://localhost:8080/todos"

    let list:{ id: number; title: string; completed: boolean }[] = []

    let newTask:string = ""

    let nextTaskId:number = 10

    async function getAllTasks() {
        fetch(URL)
        .then(response => response.json())
        .then(data => {
            console.log(data)
            list = data;
        }).catch(error => {
            console.log(error)
        })
    }

    async function handleSubmit(e: any){
        const item = {
            id: nextTaskId,
            title: newTask,
            completed: false
        }

        console.log(JSON.stringify(item))

        await fetch(URL, 
            {
                method: "POST",
                body: JSON.stringify(item),
                headers: {'Content-Type': 'application/json'},
            }
        )

        newTask = ""
        nextTaskId += 1

        getAllTasks()
    }

    async function handleDelete(id: number){
        await fetch(URL + "/" + id,
            {
                method: "DELETE"
            }
        )
        getAllTasks();
    }

    onMount(() => {
        console.log("mounted: ");
        getAllTasks()
    });

</script>

<main>

    <h1>Todos</h1>

    <div class="tasks">

    <form on:submit|preventDefault={handleSubmit}>

        <input bind:value={newTask} class="enter" type="text" placeholder="What to be done?" />

    </form>

    {#each list as t }

        <Task {handleDelete} task={t} />
      
    {/each}

  </div>

</main>

<style> 

main {
    display: flex;
    align-items: center;
    flex-direction: column;
  }
  h1 {
    color: #ccc;
    font-weight: 300;
    font-size: 8rem;
  }
  .tasks {
    width: 30rem;
    box-shadow: -5px 5px 10px -5px rgb(23 54 71 / 50%);
  }
  .enter {
    width: 100%;
    padding: 0.5rem;
    border: none;
    font-size: 1.5rem;
    outline: none;
    border-bottom: 3px solid #ddd;
  }
  .enter::placeholder { 
    color: #ccc;
    font-style: italic;
    opacity: 1;
  }

</style>
