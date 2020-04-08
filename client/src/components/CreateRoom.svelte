<script>
    import { navigate } from 'svelte-navaid';
    let roomname = "", password = "";
    
    const handleSubmission = (event) => {
        event.preventDefault();
        fetch("/createRoom", {method: "POST", body: JSON.stringify({name: roomname})})
            .then(res => {
                console.log(res);
                navigate(`#/r/${roomname}`);
            })
            .catch(error => {
                console.log(error);
                alert(error);
            });
    }
</script>

<div class="h-screen w-screen flex justify-center align-center">
    <form on:submit={handleSubmission} class="max-w-sm rounded overflow-hidden shadow-inner px-6 py-6 my-auto mx-auto h-auto ">
    <div class="grid grid-cols-1:2 grid-rows-3 gap-2">
        <label for="roomname" class="col-span-1 row-start-1">Room Name</label>
        <input name="roomname" class="col-span-1 row-start-1 shadow focus:outline-none focus:shadow-outline" autocomplete="off" type="text" bind:value={roomname} required>
        <!-- <label for="password" class="col-span-1 row-start-1">Password</label>
        <input name="password" class="col-span-1 row-start-1 shadow focus:outline-none focus:shadow-outline" autocomplete="off" type="text" bind:value={password} required> -->
        <div class="col-span-2 row-start-3 flex justify-center">
            <button type="submit" class="mx-auto rounded px-4 py-1 shadow">CREATE</button>
        </div>
    </div>
    </form>
</div>

<!-- <script>
import { Sveltik, Form, Field, ErrorMessage } from 'sveltik';

let initialValues = {
    roomname: '',
    password: '',
}

let validate = values => {
    const errors = {}
    if (!values.roomname) {
        errors.roomname = 'Required'
    } else if (
        !/^[A-Z0-9._%+-]+/i.test(values.roomname)
    ) {
        errors.roomname = 'Invalid roomname'
        console.log('invalid roomname');
    }
    return errors
}

let onSubmit = (values, { setSubmitting }) => {
    setTimeout(() => {
        alert(JSON.stringify(values, null, 2))
        setSubmitting(false)
    }, 400)
}
</script>

<div class="flex justify-center align-center h-screen">
<Sveltik {initialValues} {validate} {onSubmit} let:isSubmitting>
    <Form class="max-w-sm rounded overflow-hidden shadow-lg px-6 my-auto mx-auto py-4 h-auto ">
        <div class="grid grid-cols-1:2 grid-rows-3 gap-2">
            <label for="roomname" class="col-span-1 row-start-1">roomname</label>
            <Field type="roomname" name="roomname" class="col-span-2 row-start-1" />
            <label for="password" class="col-span-1 row-start-2">password</label>
            <Field type="password" name="password" class="col-span-1 row-start-2" />
            <button type="submit" class="col-span-2 row-start-3 border-solid border-2 border-black rounded" disabled={isSubmitting}>
                Login
            </button>
        </div>
            <ErrorMessage name="roomname" as="div" class="col-span-2 row-start-3" />
            <ErrorMessage name="password" as="div" class="col-span-2 row-start-3" />
    </Form>
</Sveltik>
</div> -->
