function main() {
    let r0 = document.getElementById("r0");
    let bytes = r0.getElementsByClassName("byte");

    let step = document.getElementById("step");
    step.addEventListener("click", () => {
        fetch("/vm/state")
        .then((response) => response.json())
        .then((state) => {
            console.log(state.regs);
        });
    });
}

main();
