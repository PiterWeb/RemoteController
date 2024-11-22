<script>
    
    import BackwardButton from '$lib/layout/BackwardButton.svelte';
    /** @type {{children?: import('svelte').Snippet}} */
    let { children } = $props();

</script>

<!-- BackwardButton for going back to config tab -->
<BackwardButton path="/mode/config"/>

{@render children?.()}