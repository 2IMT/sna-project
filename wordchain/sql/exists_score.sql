SELECT EXISTS {
    SELECT *
    FROM
        scores
    WHERE
        id = $1
};
