// example Nextflow pipeline that creates a table that we want to import into MultiQC

process DO_THING {
    // debug true
    errorStrategy "ignore"

    input:
    val(x)

    output:
    tuple val(sampleID), val("${task.process}"), topic: passed

    script:
    sampleID = "${x}"
    """
    echo "got $x"
    if [ "$x" == "Sample1" ]; then echo "bad sample!"; exit 1; fi
    """
}

process DO_THING2 {
    // debug true
    errorStrategy "ignore"

    input:
    val(x)

    output:
    tuple val(sampleID), val("${task.process}"), topic: passed

    script:
    sampleID = "${x}"
    """
    echo "got $x"
    if [ "$x" == "Sample2" ]; then echo "bad sample!"; exit 1; fi
    """
}

process DO_THING3 {
    // debug true
    errorStrategy "ignore"

    input:
    val(x)

    output:
    tuple val(sampleID), val("${task.process}"), topic: passed

    script:
    sampleID = "${x}"
    """
    echo "got $x"
    if [ "$x" == "Sample3" ]; then echo "bad sample!"; exit 1; fi
    """
}

process CONVERT_TABLE {
    publishDir "output", mode: 'copy'

    input:
    path(input_table)

    output:
    path("table.yml"), emit: table

    script:
    """
    multiqc-table-generator "$input_table" > table.yml
    """
}

process MULTIQC {
    // put your multiqc here to do thing with the table you made
    debug true
    publishDir "output", mode: 'copy'

    input:
    path(config_file)

    output:
    path("*.html")

    script:
    """
    multiqc --force --config "${config_file}" .
    """
}

workflow {
    samples = Channel.from("Sample1", "Sample2", "Sample3", "Sample4")

    // list all the input samples
    inputSamples = samples.map { sampleID ->
        return [sampleID, "INPUT"]
    }

    DO_THING(samples)
    DO_THING2(samples)
    DO_THING3(samples)

    // view the samples that passed
    allSamples = inputSamples.concat(channel.topic("passed")).map { sampleID, items ->
        // NOTE: need to force all the sampleID's back to the default string type
        // you can check the data types with; println sampleID.getClass()
        return [sampleID.toString(), items]
    }.groupTuple()
    // allSamples.view()

    // make a table out of the passing samples
    samplesTable = allSamples.map { sampleId, processList  ->
            return "${sampleId}\t${processList}"
        }
        .collectFile(name: "passed.tsv", storeDir: "output", newLine: true)
    CONVERT_TABLE(samplesTable)
    MULTIQC(CONVERT_TABLE.out.table)
}