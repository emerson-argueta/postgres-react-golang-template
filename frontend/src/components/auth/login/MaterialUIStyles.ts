import { makeStyles, Theme } from "@material-ui/core";

export const useStyles = makeStyles((theme:Theme) => ({
    paper: {
      marginTop: theme.spacing(8),
      display: 'flex',
      flexDirection: 'column',
      alignItems: 'center',
    },
    textField: {
        width: '90%',
        margin: '5%',
        marginTop:'5%'
    }
  }));