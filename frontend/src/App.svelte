<script>
    import {SyncLoader} from "svelte-loading-spinners";
    import Toastify from 'toastify-js';
    import 'toastify-js/src/toastify.css';

    const serverUrl = "http://localhost:8080";
    $: inputUrl = "";
    $: showShortenedUrl = false;
    $: shortenedUrl = "";
    $: showSpinner = false;

    function updateUrl(event) {
        inputUrl = event.currentTarget.value;
        console.log(inputUrl)
    }

    function shortenUrl(event) {
        showSpinner = true;
        console.log(showSpinner)
        event.preventDefault();

        // send the url to the api to shorten using url encoded form data
        fetch(`${serverUrl}/shorten`, {
            method: "POST",
            headers: {
                "Content-Type": "application/x-www-form-urlencoded",
            },
            body: new URLSearchParams({
                url: inputUrl,
            }),
        })
            .then((response) => {
                response.text().then((data) => {
                    shortenedUrl = `${serverUrl}/l/${data}`;
                    showShortenedUrl = true;
                    showSpinner = false;

                    // clear the input field
                    console.log(event.target.parentElement.parentElement.reset())

                    // show a toast message
                    Toastify({
                        text: "Link shortened successfully! Click here to copy",
                        duration: 4500,
                        close: false,
                        gravity: "bottom",
                        position: "center",
                        backgroundColor: "black",
                        stopOnFocus: true,
                        onClick: function(){
                            navigator.clipboard.writeText(shortenedUrl);

                            // show extremely brief "Copied!" above where the user clicked
                            Toastify({
                                text: "Copied link!",
                                duration: 1000,
                                close: false,
                                gravity: "top",
                                position: "right",
                                backgroundColor: "black",
                            }).showToast();
                        }
                    }).showToast();

                });
            })
            .catch((error) => {
                console.error("Error:", error);
            });


    }

</script>

<main>
    <div class="flex flex-col h-full w-full">
        <div class="flex justify-center font-semibold mt-5 text-[45px] w-full">
            <h1>Shine-ify URL</h1>
        </div>

        <div class="flex w-full h-full justify-center items-center">
            <form class="flex flex-col items-center justify-between border w-1/3 h-[400px] px-6 pb-4 rounded-md">
                <h2 class="w-full mt-4 mb-8 text-[18px] font-semibold">Shorten any link</h2>

                <div class="w-full h-full">
                    <div class="w-full">
                        <label class="w-full pb-[3px]" for="urlInput">URL</label>
                        <input onkeyup={updateUrl} id="urlInput" class="w-full border rounded h-[38px] pl-3" placeholder="https://google.com">
                    </div>
                </div>

                <div class="w-full">
                    {#if showShortenedUrl}
                        <div class="text-[14px] text-gray-400 mb-3">your new & shiny link - <a class="font-bold text-blue-600" href={shortenedUrl}>{shortenedUrl}</a></div>
                    {/if}
                    <button onclick={shortenUrl} class="flex justify-center items-center bg-black w-full h-[42px] rounded-lg text-white text-[16px]" type="submit">
                        {#if showSpinner}
                            <SyncLoader color="#FFFFFF" duration="1s" pause={!showSpinner} />
                        {:else}
                            Shorten
                        {/if}
                    </button>
                </div>
            </form>
        </div>
    </div>
</main>
