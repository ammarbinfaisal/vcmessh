<script>
    let first_name = "",
        last_name = "",
        father_name = "",
        mother_name = "",
        student_class = "",
        section = "";
    const handleSubmission = (event) => {
        event.preventDefault();
        const body = JSON.stringify({first_name, last_name, father_name, mother_name, class: student_class, section});
    }
</script>

<div class="wrapper">
    <form on:submit={handleSubmission}>
        <div>
            <label for="first_name"></label>
            <input name="first_name" type="text" bind:value={first_name} required>
        </div>
        <div>
            <label for="last_name"></label>
            <input name="last_name" type="text" bind:value={last_name}>
        </div>
        <div>
            <label for="father_name"></label>
            <input name="father_name" type="text" bind:value={father_name}>
        </div>
        <div>
            <label for="mother_name"></label>
            <input name="mother_name" type="text" bind:value={mother_name}>
        </div>
        <div>
            <label for="class"></label>
            <input name="class" type="text" required bind:value={student_class}>
        </div>
        <div>
            <label for="section"></label>
            <input name="section" type="text" required bind:value={section}>
        </div>
        <button type="submit"/>
    </form>
</div>
